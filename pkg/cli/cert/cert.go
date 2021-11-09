package cert

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/erikdubbelboer/gspt"
	"github.com/otiai10/copy"
	"github.com/rancher/k3s/pkg/cli/cmds"
	"github.com/rancher/k3s/pkg/daemons/config"
	"github.com/rancher/k3s/pkg/server"
	"github.com/rancher/k3s/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func commandSetup(app *cli.Context, cfg *cmds.Server, sc *server.Config) (string, error) {
	gspt.SetProcTitle(os.Args[0])

	nodeName := app.String("node-name")
	if nodeName == "" {
		h, err := os.Hostname()
		if err != nil {
			return "", err
		}
		nodeName = h
	}

	os.Setenv("NODE_NAME", nodeName)

	sc.DisableAgent = true
	sc.ControlConfig.DataDir = cfg.DataDir
	sc.ControlConfig.Runtime = &config.ControlRuntime{}

	return server.ResolveDataDir(cfg.DataDir)
}

func Rotate(app *cli.Context) error {
	if err := cmds.InitLogging(); err != nil {
		return err
	}
	return rotate(app, &cmds.ServerConfig)
}

func rotate(app *cli.Context, cfg *cmds.Server) error {
	var serverConfig server.Config

	dataDir, err := commandSetup(app, cfg, &serverConfig)
	if err != nil {
		return err
	}

	serverConfig.ControlConfig.DataDir = dataDir
	tlsDir := filepath.Join(serverConfig.ControlConfig.DataDir, "tls")
	tlsBackupDir := filepath.Join(serverConfig.ControlConfig.DataDir, "tls-"+strconv.Itoa(int(time.Now().Unix())))

	// backing up tls dir
	if _, err := os.Stat(tlsDir); err != nil {
		return err
	}
	if err := copy.Copy(tlsDir, tlsBackupDir); err != nil {
		return err
	}
	certList := []string{
		filepath.Join(dataDir, "tls", "client-admin.crt"),
		filepath.Join(dataDir, "tls", "client-admin.key"),
		filepath.Join(dataDir, "tls", "client-controller.crt"),
		filepath.Join(dataDir, "tls", "client-controller.key"),
		filepath.Join(dataDir, "tls", "client-"+version.Program+"-cloud-controller.crt"),
		filepath.Join(dataDir, "tls", "client-"+version.Program+"-cloud-controller.key"),
		filepath.Join(dataDir, "tls", "client-scheduler.crt"),
		filepath.Join(dataDir, "tls", "client-scheduler.key"),
		filepath.Join(dataDir, "tls", "client-kube-apiserver.crt"),
		filepath.Join(dataDir, "tls", "client-kube-apiserver.key"),
		filepath.Join(dataDir, "tls", "client-kube-proxy.crt"),
		filepath.Join(dataDir, "tls", "client-kube-proxy.key"),
		filepath.Join(dataDir, "tls", "client-"+version.Program+"-controller.crt"),
		filepath.Join(dataDir, "tls", "client-"+version.Program+"-controller.key"),
		filepath.Join(dataDir, "tls", "serving-kube-apiserver.crt"),
		filepath.Join(dataDir, "tls", "serving-kube-apiserver.key"),
		filepath.Join(dataDir, "tls", "client-kubelet.key"),
		filepath.Join(dataDir, "tls", "serving-kubelet.key"),
		filepath.Join(dataDir, "tls", "client-auth-proxy.crt"),
		filepath.Join(dataDir, "tls", "client-auth-proxy.key"),
		filepath.Join(dataDir, "tls", "etcd", "server-client.crt"),
		filepath.Join(dataDir, "tls", "etcd", "server-client.key"),
		filepath.Join(dataDir, "tls", "etcd", "peer-server-client.crt"),
		filepath.Join(dataDir, "tls", "etcd", "peer-server-client.key"),
		filepath.Join(dataDir, "tls", "etcd", "client.crt"),
		filepath.Join(dataDir, "tls", "etcd", "client.key"),
	}

	for _, cert := range certList {
		if err := os.Remove(cert); err != nil {
			logrus.Warnf("Certificate %s is already deleted", cert)
		}
	}
	return nil
}
