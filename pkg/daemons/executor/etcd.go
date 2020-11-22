// +build !no_embedded_executor

package executor

import (
	"github.com/rancher/k3s/pkg/util"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/embed"
	"go.etcd.io/etcd/etcdserver"
	"strings"
)

func (e Embedded) CurrentETCDOptions() (InitialOptions, error) {
	return InitialOptions{}, nil
}

func (e Embedded) ETCD(args ETCDConfig) error {
	configFile, err := args.ToConfigFile()
	if err != nil {
		return err
	}
	cfg, err := embed.ConfigFromFile(configFile)
	if err != nil {
		return err
	}
	etcd, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil
	}

	go func() {
		select {
		case err := <-etcd.Server.ErrNotify():
			var backupdatadir string
			if strings.Contains(err.Error(), etcdserver.ErrMemberRemoved.Error()) {
				if backupdatadir, err = util.BackupDirWithRetention(args.DataDir, 5); err != nil {
					logrus.Fatalf("Failed to remove old etcd datadir: %v", err)
				}
			}
			logrus.Fatalf("etcd data dir was moved to %s - please re run k3s and it will join the cluster", backupdatadir)
		case <-etcd.Server.StopNotify():
			logrus.Fatalf("etcd stopped")
		case err := <-etcd.Err():
			logrus.Fatalf("etcd exited: %v", err)
		}
	}()
	return nil
}