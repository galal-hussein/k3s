package cluster

// A managed database is one whose lifecycle we control - initializing the cluster, adding/removing members, taking snapshots, etc.
// This is currently just used for the embedded etcd datastore. Kine and other external etcd clusters are NOT considered managed.

import (
	"context"
	"fmt"
	"github.com/rancher/k3s/pkg/version"
	"net/http"
	"strings"
	"os"
	"path/filepath"
	"time"

	"github.com/rancher/k3s/pkg/cluster/managed"
	"github.com/rancher/kine/pkg/endpoint"
	"github.com/sirupsen/logrus"
)

// testClusterDB returns a channel that will be closed when the datastore connection is available.
// The datastore is tested for readiness every 5 seconds until the test succeeds.
func (c *Cluster) testClusterDB(ctx context.Context) (<-chan struct{}, error) {
	result := make(chan struct{})
	if c.managedDB == nil {
		close(result)
		return result, nil
	}

	go func() {
		defer close(result)
		for {
			if err := c.managedDB.Test(ctx, c.clientAccessInfo); err != nil {
				logrus.Infof("Failed to test data store connection: %v", err)
			} else {
				logrus.Info(c.managedDB.EndpointName() + " data store connection OK")
				return
			}

			select {
			case <-time.After(5 * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}()

	return result, nil
}

// start starts the database, unless a cluster reset has been requested, in which case
// it does that instead.
func (c *Cluster) start(ctx context.Context) error {
	resetFile := filepath.Join(c.config.DataDir, "db", "reset-flag")
	fmt.Println("hhhhhhhhhhhh", c.managedDB)
	if c.managedDB == nil {
		return nil
	}
	fmt.Println("hgalal.........etcd managed 1")

	if c.config.ClusterReset {
		if _, err := os.Stat(resetFile); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("cluster-reset was successfully performed, " +
					"please remove the cluster-reset flag and start %s normally, " +
					"if you need to perform another cluster reset, " +
					"you must first manually delete the %s file",
					version.Program, resetFile)
			}
		}
		return c.managedDB.Reset(ctx, c.clientAccessInfo)
	}
	fmt.Println("hgalal.........etcd managed")
	// removing the reset file and ignore error if the file doesnt exist
	os.Remove(resetFile)

	return c.managedDB.Start(ctx, c.clientAccessInfo)
}

// initClusterDB registers routes for database info with the http request handler
func (c *Cluster) initClusterDB(ctx context.Context, handler http.Handler) (http.Handler, error) {
	if c.managedDB == nil {
		return handler, nil
	}

	if !strings.HasPrefix(c.config.Datastore.Endpoint, c.managedDB.EndpointName()+"://") {
		c.config.Datastore = endpoint.Config{
			Endpoint: c.managedDB.EndpointName(),
		}
	}

	return c.managedDB.Register(ctx, c.config, handler)
}

// assignManagedDriver assigns a driver based on a number of different configuration variables.
// If a driver has been initialized it is used.
// If the configured endpoint matches the name of a driver, that driver is used.
// If no specific endpoint has been requested and creating or joining has been requested,
// we use the default driver.
// If none of the above are true, no managed driver is assigned.
func (c *Cluster) assignManagedDriver(ctx context.Context) error {
	// Check all managed drivers for an initialized database on disk; use one if found
	for _, driver := range managed.Registered() {
		if ok, err := driver.IsInitialized(ctx, c.config); err != nil {
			return err
		} else if ok {
			c.managedDB = driver
			return nil
		}
	}

	// This is needed to allow downstreams to override driver selection logic by
	// setting ServerConfig.Datastore.Endpoint such that it will match a driver's EndpointName
	endpointType := strings.SplitN(c.config.Datastore.Endpoint, ":", 2)[0]
	for _, driver := range managed.Registered() {
		if endpointType == driver.EndpointName() {
			c.managedDB = driver
			return nil
		}
	}

	// If we have been asked to initialize or join a cluster, do so using the default managed database.
	if c.config.Datastore.Endpoint == "" && (c.config.ClusterInit || (c.config.Token != "" && c.config.JoinURL != "")) {
		for _, driver := range managed.Registered() {
			if driver.EndpointName() == managed.Default() {
				c.managedDB = driver
				return nil
			}
		}
	}

	return nil
}
