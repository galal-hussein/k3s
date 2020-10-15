package etcd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	certutil "github.com/rancher/dynamiclistener/cert"
	"github.com/rancher/k3s/pkg/clientaccess"
	"github.com/rancher/k3s/pkg/daemons/config"
	"github.com/rancher/k3s/pkg/daemons/executor"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	etcd "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/snapshot"
	"go.etcd.io/etcd/etcdserver/etcdserverpb"
	utilnet "k8s.io/apimachinery/pkg/util/net"
)

type ETCD struct {
	client  *etcd.Client
	config  *config.Control
	name    string
	runtime *config.ControlRuntime
	address string
	cron    *cron.Cron
}

// NewETCD creates a new value of type
// ETCD with an initialized cron value.
func NewETCD() *ETCD {
	return &ETCD{
		cron: cron.New(),
	}
}

const (
	snapshotPrefix = "etcd-snapshot-"
	endpoint       = "https://127.0.0.1:2379"
	testTimeout    = time.Second * 10

	// defaults from etcdctl/ctlv3/ctl.go
	defaultDialTimeout      = 2 * time.Second
	defaultKeepAliveTime    = 2 * time.Second
	defaultKeepAliveTimeOut = 6 * time.Second
)

// Members contains a slice that holds all
// members of the cluster.
type Members struct {
	Members []*etcdserverpb.Member `json:"members"`
}

// EndpointName returns the name of the endpoint.
func (e *ETCD) EndpointName() string {
	return "etcd"
}

// Test ensures that the local node is a part of the target cluster. If it is a learner, a goroutine
// will be started to promote it to full member. If it is not a part of the cluster, an error is raised.
func (e *ETCD) Test(ctx context.Context, clientAccessInfo *clientaccess.Info) error {
	ctx, cancel := context.WithTimeout(ctx, testTimeout)
	defer cancel()
	status, err := e.client.Status(ctx, endpoint)
	if err != nil {
		return err
	}

	if status.IsLearner {
		if err := e.promoteMember(ctx, clientAccessInfo); err != nil {
			return err
		}
	}
	members, err := e.client.MemberList(ctx)
	if err != nil {
		return err
	}

	var memberNameUrls []string
	for _, member := range members.Members {
		for _, peerURL := range member.PeerURLs {
			if peerURL == e.peerURL() && e.name == member.Name {
				return nil
			}
		}
		if len(member.PeerURLs) > 0 {
			memberNameUrls = append(memberNameUrls, member.Name+"="+member.PeerURLs[0])
		}
	}
	msg := fmt.Sprintf("This server is a not a member of the etcd cluster. Found %v, expect: %s=%s", memberNameUrls, e.name, e.address)
	logrus.Error(msg)
	return fmt.Errorf(msg)
}

// etcdDBDir returns the path to dataDir/db/etcd
func etcdDBDir(config *config.Control) string {
	return filepath.Join(config.DataDir, "db", "etcd")
}

// walDir returns the path to etcdDBDir/member/wal
func walDir(config *config.Control) string {
	return filepath.Join(etcdDBDir(config), "member", "wal")
}

// nameFile returns the path to etcdDBDir/name
func nameFile(config *config.Control) string {
	return filepath.Join(etcdDBDir(config), "name")
}

// ResetFile returns the path to etcdDBDir/reset-flag
func ResetFile(config *config.Control) string {
	return filepath.Join(config.DataDir, "db", "reset-flag")
}

// IsInitialized checks to see if a WAL directory exists. If so, we assume that etcd
// has already been brought up at least once.
func (e *ETCD) IsInitialized(ctx context.Context, config *config.Control) (bool, error) {
	dir := walDir(config)
	if s, err := os.Stat(dir); err == nil && s.IsDir() {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.Wrapf(err, "invalid state for wal directory %s", dir)
	}
}

// Reset resets an etcd node
func (e *ETCD) Reset(ctx context.Context, clientAccessInfo *clientaccess.Info) error {
	// Wait for etcd to come up as a new single-node cluster, then exit
	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for range t.C {
			if err := e.Test(ctx, clientAccessInfo); err == nil {
				members, err := e.client.MemberList(ctx)
				if err != nil {
					continue
				}

				if len(members.Members) == 1 && members.Members[0].Name == e.name {
					logrus.Infof("Etcd is running, restart without --cluster-reset flag now. Backup and delete ${datadir}/server/db on each peer etcd server and rejoin the nodes")
					os.Exit(0)
				}
			}
		}
	}()

	// If asked to restore from a snapshot, do so
	if e.config.ClusterResetRestorePath != "" {
		info, err := os.Stat(e.config.ClusterResetRestorePath)
		if os.IsNotExist(err) {
			return fmt.Errorf("etcd: snapshot path does not exist: %s", e.config.ClusterResetRestorePath)
		}
		if info.IsDir() {
			return fmt.Errorf("etcd: snapshot path must be a file, not a directory: %s", e.config.ClusterResetRestorePath)
		}
		if err := e.Restore(ctx); err != nil {
			return err
		}
	}

	if err := e.setName(true); err != nil {
		return err
	}
	// touch a file to avoid multiple resets
	if err := ioutil.WriteFile(ResetFile(e.config), []byte{}, 0600); err != nil {
		return err
	}
	return e.newCluster(ctx, true)
}

// Start starts the datastore
func (e *ETCD) Start(ctx context.Context, clientAccessInfo *clientaccess.Info) error {
	existingCluster, err := e.IsInitialized(ctx, e.config)
	if err != nil {
		return errors.Wrapf(err, "configuration validation failed")
	}

	e.config.Runtime.ClusterControllerStart = func(ctx context.Context) error {
		Register(ctx, e, e.config.Runtime.Core.Core().V1().Node())
		return nil
	}

	if !e.config.EtcdDisableSnapshots {
		e.setSnapshotFunction(ctx)
		e.cron.Start()
	}

	if existingCluster {
		opt, err := executor.CurrentETCDOptions()
		if err != nil {
			return err
		}
		return e.cluster(ctx, false, opt)
	}

	if clientAccessInfo == nil {
		return e.newCluster(ctx, false)
	}
	err = e.join(ctx, clientAccessInfo)
	return errors.Wrap(err, "joining etcd cluster")
}

// join attempts to add a member to an existing cluster
func (e *ETCD) join(ctx context.Context, clientAccessInfo *clientaccess.Info) error {
	clientURLs, memberList, err := e.clientURLs(ctx, clientAccessInfo)
	if err != nil {
		return err
	}

	client, err := getClient(ctx, e.runtime, clientURLs...)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var (
		cluster []string
		add     = true
	)

	members, err := client.MemberList(ctx)
	if err != nil {
		logrus.Errorf("Failed to get member list from etcd cluster. Will assume this member is already added")
		members = &etcd.MemberListResponse{
			Members: append(memberList.Members, &etcdserverpb.Member{
				Name:     e.name,
				PeerURLs: []string{e.peerURL()},
			}),
		}
		add = false
	}

	for _, member := range members.Members {
		for _, peer := range member.PeerURLs {
			u, err := url.Parse(peer)
			if err != nil {
				return err
			}
			// An uninitialized member won't have a name
			if u.Hostname() == e.address && (member.Name == e.name || member.Name == "") {
				add = false
			}
			if member.Name == "" && u.Hostname() == e.address {
				member.Name = e.name
			}
			if len(member.PeerURLs) > 0 {
				cluster = append(cluster, fmt.Sprintf("%s=%s", member.Name, member.PeerURLs[0]))
			}
		}
	}

	if add {
		logrus.Infof("Adding %s to etcd cluster %v", e.peerURL(), cluster)
		if _, err = client.MemberAddAsLearner(ctx, []string{e.peerURL()}); err != nil {
			return err
		}
		cluster = append(cluster, fmt.Sprintf("%s=%s", e.name, e.peerURL()))
	}

	go e.promoteMember(ctx, clientAccessInfo)

	logrus.Infof("Starting etcd for cluster %v", cluster)
	return e.cluster(ctx, false, executor.InitialOptions{
		Cluster: strings.Join(cluster, ","),
		State:   "existing",
	})
}

// Register configures a new etcd client and adds db info routes for the http request handler.
func (e *ETCD) Register(ctx context.Context, config *config.Control, handler http.Handler) (http.Handler, error) {
	e.config = config
	e.runtime = config.Runtime

	client, err := getClient(ctx, e.runtime, endpoint)
	if err != nil {
		return nil, err
	}
	e.client = client

	address, err := getAdvertiseAddress(config.AdvertiseIP)
	if err != nil {
		return nil, err
	}
	e.address = address

	e.config.Datastore.Endpoint = endpoint
	e.config.Datastore.Config.CAFile = e.runtime.ETCDServerCA
	e.config.Datastore.Config.CertFile = e.runtime.ClientETCDCert
	e.config.Datastore.Config.KeyFile = e.runtime.ClientETCDKey

	if err := e.setName(false); err != nil {
		return nil, err
	}

	return e.handler(handler), err
}

// setName sets a unique name for this cluster member. The first time this is called,
// or if force is set to true, a new name will be generated and written to disk. The persistent
// name is used on subsequent calls.
func (e *ETCD) setName(force bool) error {
	fileName := nameFile(e.config)
	data, err := ioutil.ReadFile(fileName)
	if os.IsNotExist(err) || force {
		h, err := os.Hostname()
		if err != nil {
			return err
		}
		e.name = strings.SplitN(h, ".", 2)[0] + "-" + uuid.New().String()[:8]
		if err := os.MkdirAll(filepath.Dir(fileName), 0700); err != nil {
			return err
		}
		return ioutil.WriteFile(fileName, []byte(e.name), 0600)
	} else if err != nil {
		return err
	}
	e.name = string(data)
	return nil
}

// handler wraps the handler with routes for database info
func (e *ETCD) handler(next http.Handler) http.Handler {
	mux := mux.NewRouter()
	mux.Handle("/db/info", e.infoHandler())
	mux.NotFoundHandler = next
	return mux
}

// infoHandler returns etcd cluster information. This is used by new members when joining the custer.
func (e *ETCD) infoHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()

		members, err := e.client.MemberList(ctx)
		if err != nil {
			json.NewEncoder(rw).Encode(&Members{
				Members: []*etcdserverpb.Member{
					{
						Name:       e.name,
						PeerURLs:   []string{e.peerURL()},
						ClientURLs: []string{e.clientURL()},
					},
				},
			})
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&Members{
			Members: members.Members,
		})
	})
}

// getClient returns an etcd client connected to the specified endpoints
func getClient(ctx context.Context, runtime *config.ControlRuntime, endpoints ...string) (*etcd.Client, error) {
	cfg, err := getClientConfig(ctx, runtime, endpoints...)
	if err != nil {
		return nil, err
	}
	return etcd.New(*cfg)
}

//getClientConfig generates an etcd client config connected to the specified endpoints
func getClientConfig(ctx context.Context, runtime *config.ControlRuntime, endpoints ...string) (*etcd.Config, error) {
	tlsConfig, err := toTLSConfig(runtime)
	if err != nil {
		return nil, err
	}

	cfg := &etcd.Config{
		Endpoints:            endpoints,
		TLS:                  tlsConfig,
		Context:              ctx,
		DialTimeout:          defaultDialTimeout,
		DialKeepAliveTime:    defaultKeepAliveTime,
		DialKeepAliveTimeout: defaultKeepAliveTimeOut,
	}

	return cfg, nil
}

// toTLSConfig converts the ControlRuntime configuration to TLS configuration suitable
// for use by etcd.
func toTLSConfig(runtime *config.ControlRuntime) (*tls.Config, error) {
	clientCert, err := tls.LoadX509KeyPair(runtime.ClientETCDCert, runtime.ClientETCDKey)
	if err != nil {
		return nil, err
	}

	pool, err := certutil.NewPool(runtime.ETCDServerCA)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{clientCert},
	}, nil
}

// getAdvertiseAddress returns the IP address best suited for advertising to clients
func getAdvertiseAddress(advertiseIP string) (string, error) {
	ip := advertiseIP
	if ip == "" {
		ipAddr, err := utilnet.ChooseHostInterface()
		if err != nil {
			return "", err
		}
		ip = ipAddr.String()
	}

	return ip, nil
}

// newCluster returns options to set up etcd for a new cluster
func (e *ETCD) newCluster(ctx context.Context, reset bool) error {
	return e.cluster(ctx, reset, executor.InitialOptions{
		AdvertisePeerURL: fmt.Sprintf("https://%s:2380", e.address),
		Cluster:          fmt.Sprintf("%s=https://%s:2380", e.name, e.address),
		State:            "new",
	})
}

// peerURL returns the peer access address for the local node
func (e *ETCD) peerURL() string {
	return fmt.Sprintf("https://%s:2380", e.address)
}

// clientURL returns the client access address for the local node
func (e *ETCD) clientURL() string {
	return fmt.Sprintf("https://%s:2379", e.address)
}

// cluster returns ETCDConfig for a cluster
func (e *ETCD) cluster(ctx context.Context, forceNew bool, options executor.InitialOptions) error {
	return executor.ETCD(executor.ETCDConfig{
		Name:                e.name,
		InitialOptions:      options,
		ForceNewCluster:     forceNew,
		ListenClientURLs:    fmt.Sprintf(e.clientURL() + ",https://127.0.0.1:2379"),
		ListenMetricsURLs:   "http://127.0.0.1:2381",
		ListenPeerURLs:      e.peerURL(),
		AdvertiseClientURLs: e.clientURL(),
		DataDir:             etcdDBDir(e.config),
		ServerTrust: executor.ServerTrust{
			CertFile:       e.config.Runtime.ServerETCDCert,
			KeyFile:        e.config.Runtime.ServerETCDKey,
			ClientCertAuth: true,
			TrustedCAFile:  e.config.Runtime.ETCDServerCA,
		},
		PeerTrust: executor.PeerTrust{
			CertFile:       e.config.Runtime.PeerServerClientETCDCert,
			KeyFile:        e.config.Runtime.PeerServerClientETCDKey,
			ClientCertAuth: true,
			TrustedCAFile:  e.config.Runtime.ETCDPeerCA,
		},
		ElectionTimeout:   5000,
		HeartbeatInterval: 500,
	})
}

// removePeer removes a peer from the cluster. The peer ID and IP address must both match.
func (e *ETCD) removePeer(ctx context.Context, id, address string) error {
	members, err := e.client.MemberList(ctx)
	if err != nil {
		return err
	}

	for _, member := range members.Members {
		if member.Name != id {
			continue
		}
		for _, peerURL := range member.PeerURLs {
			u, err := url.Parse(peerURL)
			if err != nil {
				return err
			}
			if u.Hostname() == address {
				logrus.Infof("Removing name=%s id=%d address=%s from etcd", member.Name, member.ID, address)
				_, err := e.client.MemberRemove(ctx, member.ID)
				if strings.Contains(err.Error(), rpctypes.ErrGRPCMemberNotFound.Error()) {
					logrus.Infof("husssssss: member not found returning nil")
					return nil
				}
				return err
			}
		}
	}

	return nil
}

// promoteMember attempts to promote any learners to full members at 5 second intervals.
// It will return when a member has been promoted. Usually this function is run on the node
// that has just been added to the cluster and is trying to promote itself. If it is run when there
// are no learners, it will never return.
func (e *ETCD) promoteMember(ctx context.Context, clientAccessInfo *clientaccess.Info) error {
	clientURLs, _, err := e.clientURLs(ctx, clientAccessInfo)
	if err != nil {
		return err
	}
	memberPromoted := true
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for range t.C {
		client, err := getClient(ctx, e.runtime, clientURLs...)
		// continue on errors to keep trying to promote member
		// grpc error are shown so no need to re log them
		if err != nil {
			continue
		}
		members, err := client.MemberList(ctx)
		if err != nil {
			continue
		}
		for _, member := range members.Members {
			// only one learner can exist in the cluster
			if !member.IsLearner {
				continue
			}
			if _, err := client.MemberPromote(ctx, member.ID); err != nil {
				memberPromoted = false
				break
			}
		}
		if memberPromoted {
			break
		}
	}
	return nil
}

// clientURLs returns a list of all non-learner etcd cluster member client access URLs
func (e *ETCD) clientURLs(ctx context.Context, clientAccessInfo *clientaccess.Info) ([]string, Members, error) {
	var memberList Members
	resp, err := clientaccess.Get("/db/info", clientAccessInfo)
	if err != nil {
		return nil, memberList, err
	}

	if err := json.Unmarshal(resp, &memberList); err != nil {
		return nil, memberList, err
	}

	var clientURLs []string
	for _, member := range memberList.Members {
		// excluding learner member from the client list
		if member.IsLearner {
			continue
		}
		clientURLs = append(clientURLs, member.ClientURLs...)
	}
	return clientURLs, memberList, nil
}

// snapshotDir ensures that the snapshot directory exists, and then returns its path.
func snapshotDir(config *config.Control) (string, error) {
	if config.EtcdSnapshotDir == "" {
		// we have to create the snapshot dir if we are using
		// the default snapshot dir if it doesn't exist
		defaultSnapshotDir := filepath.Join(config.DataDir, "db", "snapshots")
		s, err := os.Stat(defaultSnapshotDir)
		if err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(defaultSnapshotDir, 0700); err != nil {
					return "", err
				}
				return defaultSnapshotDir, nil
			}
			return "", err
		}
		if s.IsDir() {
			return defaultSnapshotDir, nil
		}
	}
	return config.EtcdSnapshotDir, nil
}

// snapshot attempts to save a new snapshot to the configured directory, and then clean up any old
// snapshots in excess of the retention limits.
func (e *ETCD) snapshot(ctx context.Context) {
	status, err := e.client.Status(ctx, endpoint)
	if err != nil {
		logrus.Errorf("Failed to check etcd status for snapshot: %v", err)
		return
	}

	if status.IsLearner {
		logrus.Warnf("Skipping snapshot: not supported for learner")
		return
	}

	snapshotDir, err := snapshotDir(e.config)
	if err != nil {
		logrus.Errorf("Failed to get the snapshot dir: %v", err)
		return
	}

	cfg, err := getClientConfig(ctx, e.runtime, endpoint)
	if err != nil {
		logrus.Errorf("Failed to get config for etcd snapshot: %v", err)
		return
	}

	snapshotPath := filepath.Join(snapshotDir, snapshotPrefix+strconv.Itoa(int(time.Now().Unix())))
	logrus.Infof("Saving etcd snapshot to %s", snapshotPath)

	if err := snapshot.NewV3(nil).Save(ctx, *cfg, snapshotPath); err != nil {
		logrus.Errorf("Failed to save snapshot: %v", err)
		return
	}
	if err := snapshotRetention(e.config.EtcdSnapshotRetention, snapshotDir); err != nil {
		logrus.Errorf("Failed to apply snapshot retention: %v", err)
		return
	}
}

// setSnapshotFunction schedules snapshots at the configured interval
func (e *ETCD) setSnapshotFunction(ctx context.Context) {
	e.cron.AddFunc(e.config.EtcdSnapshotCron, func() { e.snapshot(ctx) })
}

// Restore performs a restore of the ETCD datastore from
// the given snapshot path. This operation exists upon
// completion.
func (e *ETCD) Restore(ctx context.Context) error {
	// check the old etcd data dir
	oldDataDir := etcdDBDir(e.config) + "-old-" + strconv.Itoa(int(time.Now().Unix()))
	if e.config.ClusterResetRestorePath == "" {
		return errors.New("no etcd restore path was specified")
	}
	// make sure snapshot exists before restoration
	if _, err := os.Stat(e.config.ClusterResetRestorePath); err != nil {
		return err
	}
	// move the data directory to a temp path
	if err := os.Rename(etcdDBDir(e.config), oldDataDir); err != nil {
		return err
	}
	logrus.Infof("Pre-restore etcd database moved to %s", oldDataDir)
	sManager := snapshot.NewV3(nil)
	if err := sManager.Restore(snapshot.RestoreConfig{
		SnapshotPath:   e.config.ClusterResetRestorePath,
		Name:           e.name,
		OutputDataDir:  etcdDBDir(e.config),
		OutputWALDir:   walDir(e.config),
		PeerURLs:       []string{e.peerURL()},
		InitialCluster: e.name + "=" + e.peerURL(),
	}); err != nil {
		return err
	}
	return nil
}

// snapshotRetention iterates through the snapshots and removes the oldest
// leaving the desired number of snapshots.
func snapshotRetention(retention int, snapshotDir string) error {
	var snapshotFiles []os.FileInfo
	if err := filepath.Walk(snapshotDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), snapshotPrefix) {
			snapshotFiles = append(snapshotFiles, info)
		}
		return nil
	}); err != nil {
		return err
	}
	if len(snapshotFiles) <= retention {
		return nil
	}
	sort.Slice(snapshotFiles, func(i, j int) bool {
		return snapshotFiles[i].Name() < snapshotFiles[j].Name()
	})
	return os.Remove(filepath.Join(snapshotDir, snapshotFiles[0].Name()))
}
