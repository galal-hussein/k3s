package server

import (
	"context"
	"fmt"
	net2 "net"
	"os"
	"path/filepath"
	"strings"

	systemd "github.com/coreos/go-systemd/daemon"
	"github.com/erikdubbelboer/gspt"
	"github.com/pkg/errors"
	"github.com/rancher/k3s/pkg/agent"
	"github.com/rancher/k3s/pkg/cli/cmds"
	"github.com/rancher/k3s/pkg/datadir"
	"github.com/rancher/k3s/pkg/netutil"
	"github.com/rancher/k3s/pkg/server"
	"github.com/rancher/k3s/pkg/token"
	"github.com/rancher/k3s/pkg/version"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/util/net"
	kubeapiserverflag "k8s.io/component-base/cli/flag"
	"k8s.io/kubernetes/pkg/master"

	_ "github.com/go-sql-driver/mysql" // ensure we have mysql
	_ "github.com/lib/pq"              // ensure we have postgres
	_ "github.com/mattn/go-sqlite3"    // ensure we have sqlite
)

func Run(app *cli.Context) error {
	if err := cmds.InitLogging(); err != nil {
		return err
	}
	return run(app, &cmds.ServerConfig)
}

func run(app *cli.Context, cfg *cmds.Server) error {
	var (
		err error
	)

	// hide process arguments from ps output, since they may contain
	// database credentials or other secrets.
	gspt.SetProcTitle(os.Args[0] + " etcd")

	if os.Getuid() != 0 {
		return fmt.Errorf("must run as root unless --disable-agent is specified")
	}

	if cfg.Token == "" && cfg.ClusterSecret != "" {
		cfg.Token = cfg.ClusterSecret
	}

	etcdconfig := server.Config{}
	etcdconfig.ControlConfig.Token = cfg.Token
	etcdconfig.ControlConfig.AgentToken = cfg.AgentToken
	etcdconfig.ControlConfig.JoinURL = cfg.ServerURL
	if cfg.AgentTokenFile != "" {
		etcdconfig.ControlConfig.AgentToken, err = token.ReadFile(cfg.AgentTokenFile)
		if err != nil {
			return err
		}
	}
	if cfg.TokenFile != "" {
		etcdconfig.ControlConfig.Token, err = token.ReadFile(cfg.TokenFile)
		if err != nil {
			return err
		}
	}
	etcdconfig.ControlConfig.DataDir = cfg.DataDir
	etcdconfig.ControlConfig.KubeConfigOutput = cfg.KubeConfigOutput
	etcdconfig.ControlConfig.KubeConfigMode = cfg.KubeConfigMode
	etcdconfig.ControlConfig.SANs = knownIPs(cfg.TLSSan)
	etcdconfig.ControlConfig.SupervisorPort = cfg.SupervisorPort
	etcdconfig.ControlConfig.FlannelBackend = cfg.FlannelBackend
	etcdconfig.ControlConfig.DisableNPC = cfg.DisableNPC
	etcdconfig.ControlConfig.DisableKubeProxy = cfg.DisableKubeProxy
	etcdconfig.ControlConfig.EtcdSnapshotCron = cfg.EtcdSnapshotCron
	etcdconfig.ControlConfig.EtcdSnapshotDir = cfg.EtcdSnapshotDir
	etcdconfig.ControlConfig.EtcdSnapshotRetention = cfg.EtcdSnapshotRetention
	etcdconfig.ControlConfig.EtcdDisableSnapshots = cfg.EtcdDisableSnapshots

	if cfg.ClusterResetRestorePath != "" && !cfg.ClusterReset {
		return errors.New("Invalid flag use. --cluster-reset required with --cluster-reset-restore-path")
	}

	etcdconfig.ControlConfig.ClusterReset = cfg.ClusterReset
	etcdconfig.ControlConfig.ClusterResetRestorePath = cfg.ClusterResetRestorePath

	if etcdconfig.ControlConfig.SupervisorPort == 0 {
		etcdconfig.ControlConfig.SupervisorPort = etcdconfig.ControlConfig.HTTPSPort
	}

	if cmds.AgentConfig.FlannelIface != "" && cmds.AgentConfig.NodeIP == "" {
		cmds.AgentConfig.NodeIP = netutil.GetIPFromInterface(cmds.AgentConfig.FlannelIface)
	}
	if etcdconfig.ControlConfig.PrivateIP == "" && cmds.AgentConfig.NodeIP != "" {
		etcdconfig.ControlConfig.PrivateIP = cmds.AgentConfig.NodeIP
	}

	etcdconfig.ControlConfig.SANs = append(etcdconfig.ControlConfig.SANs, apiServerServiceIP.String())

	etcdconfig.StartupHooks = append(etcdconfig.StartupHooks, cfg.StartupHooks...)

	logrus.Info("Starting " + version.Program + " " + app.App.Version)
	notifySocket := os.Getenv("NOTIFY_SOCKET")
	os.Unsetenv("NOTIFY_SOCKET")

	ctx := signals.SetupSignalHandler(context.Background())
	if err := server.StartServer(ctx, &etcdconfig); err != nil {
		return err
	}

	go func() {
		<-etcdconfig.ControlConfig.Runtime.ETCDReady
		logrus.Info("ETCD server is now running")
		logrus.Info(version.Program + " is up and running")
		if notifySocket != "" {
			os.Setenv("NOTIFY_SOCKET", notifySocket)
			systemd.SdNotify(true, "READY=1\n")
		}
	}()

	if cfg.DisableAgent {
		<-ctx.Done()
		return nil
	}

	ip := etcdconfig.ControlConfig.BindAddress
	if ip == "" {
		ip = "127.0.0.1"
	}

	url := fmt.Sprintf("https://%s:%d", ip, etcdconfig.ControlConfig.SupervisorPort)
	token, err := server.FormatToken(etcdconfig.ControlConfig.Runtime.AgentToken, etcdconfig.ControlConfig.Runtime.ServerCA)
	if err != nil {
		return err
	}

	agentConfig := cmds.AgentConfig
	agentConfig.Debug = app.GlobalBool("debug")
	agentConfig.DataDir = filepath.Dir(etcdconfig.ControlConfig.DataDir)
	agentConfig.ServerURL = url
	agentConfig.Token = token
	agentConfig.DisableLoadBalancer = true

	return agent.Run(ctx, agentConfig)
}

func knownIPs(ips []string) []string {
	ips = append(ips, "127.0.0.1")
	ip, err := net.ChooseHostInterface()
	if err == nil {
		ips = append(ips, ip.String())
	}
	return ips
}
