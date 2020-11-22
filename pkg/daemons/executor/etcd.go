// +build !no_embedded_executor

package executor

import (
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/embed"
	"go.etcd.io/etcd/etcdserver"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
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
			logrus.Info(etcdserver.ErrMemberRemoved)
			if strings.Contains(err.Error(), etcdserver.ErrMemberRemoved.Error()) {
				if backupdatadir, err = RemoveDataDir(args.DataDir); err != nil {
					logrus.Fatalf("Failed to remove old etcd datadir: %v", err)
				}
			}
			logrus.Infof("here2")
			logrus.Fatalf("etcd data dir was moved to %s - please re run k3s and it will join the cluster", backupdatadir)
		case <-etcd.Server.StopNotify():
			logrus.Fatalf("etcd stopped")
		case err := <-etcd.Err():
			logrus.Fatalf("etcd exited: %v", err)
		}
	}()
	return nil
}

// RemoveDataDir will move the datadir to a backup dir
// and will keep only maxBackupRetention of datadirs
func RemoveDataDir(datadir string) (string, error) {
	logrus.Infof("in remove data dir")
	backupDataDir := datadir + "-backup-" + strconv.Itoa(int(time.Now().Unix()))
	if _, err := os.Stat(datadir); err != nil {
		return "", nil
	}
	files, err := ioutil.ReadDir(filepath.Dir(datadir))
	if err != nil {
		return "", err
	}
	sort.Slice(files, func(i,j int) bool{
		return files[i].ModTime().After(files[j].ModTime())
	})
	count := 0
	for _, f := range files {
		logrus.Infof("here")
		count++
		logrus.Info(f)
		logrus.Info(f.IsDir())
		if strings.HasPrefix(f.Name(), "etcd-backup") && f.IsDir() {
			if count > 5 {
				logrus.Infof("deleting")
				if err := os.RemoveAll(filepath.Join(filepath.Dir(datadir), f.Name())); err != nil {
					return "", err
				}
			}
		}
	}
	// move the data directory to a temp path
	logrus.Info(backupDataDir, datadir)
	if err := os.Rename(datadir, backupDataDir); err != nil {
		return "", err
	}
	return backupDataDir, nil
}