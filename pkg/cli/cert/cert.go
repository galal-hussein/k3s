package cert

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erikdubbelboer/gspt"
	"github.com/otiai10/copy"
	"github.com/rancher/k3s/pkg/cli/cmds"
	"github.com/rancher/k3s/pkg/daemons/config"
	"github.com/rancher/k3s/pkg/datadir"
	"github.com/rancher/k3s/pkg/server"
	"github.com/rancher/k3s/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	adminComponent             = "admin"
	apiServerComponent         = "api-server"
	controllerManagerComponent = "controller-manager"
	schedulerComponent         = "scheduler"
	etcdComponent              = "etcd"
	programControllerComponent = "-controller"
	authProxyComponent         = "auth-proxy"
	cloudControllerComponent   = "cloud-controller"
	kubeletComponent           = "kubelet"
	kubeProxyComponent         = "kube-proxy"
)

func commandSetup(app *cli.Context, cfg *cmds.Server, sc *server.Config) (string, string, error) {
	gspt.SetProcTitle(os.Args[0])

	nodeName := app.String("node-name")
	if nodeName == "" {
		h, err := os.Hostname()
		if err != nil {
			return "", "", err
		}
		nodeName = h
	}

	os.Setenv("NODE_NAME", nodeName)

	sc.ControlConfig.DataDir = cfg.DataDir
	sc.ControlConfig.Runtime = &config.ControlRuntime{}
	dataDir, err := datadir.Resolve(cfg.DataDir)
	if err != nil {
		return "", "", err
	}
	return filepath.Join(dataDir, "server"), filepath.Join(dataDir, "agent"), err
}

func Run(app *cli.Context) error {
	if err := cmds.InitLogging(); err != nil {
		return err
	}
	return rotate(app, &cmds.ServerConfig)
}

func rotate(app *cli.Context, cfg *cmds.Server) error {
	var serverConfig server.Config

	serverDataDir, agentDataDir, err := commandSetup(app, cfg, &serverConfig)
	if err != nil {
		return err
	}

	serverConfig.ControlConfig.DataDir = serverDataDir
	tlsDir := filepath.Join(serverConfig.ControlConfig.DataDir, "tls")
	tlsBackupDir := filepath.Join(serverConfig.ControlConfig.DataDir, "tls-"+strconv.Itoa(int(time.Now().Unix())))

	// backing up tls dir
	if _, err := os.Stat(tlsDir); err != nil {
		return err
	}
	if err := copy.Copy(tlsDir, tlsBackupDir); err != nil {
		return err
	}
	if len(cmds.ComponentList) == 0 {
		// rotate all certs
		return rotateAllCerts(filepath.Join(serverDataDir, "tls"), agentDataDir)
	}
	certList := []string{}
	for _, component := range cmds.ComponentList {
		switch component {
		case adminComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-admin.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-admin.key"))
		case apiServerComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-kube-apiserver.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-kube-apiserver.key"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "serving-kube-apiserver.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "serving-kube-apiserver.key"))
		case controllerManagerComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-controller.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-controller.key"))
		case schedulerComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-scheduler.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-scheduler.key"))
		case etcdComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "etcd", "server-client.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "etcd", "server-client.key"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "etcd", "peer-server-client.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "etcd", "peer-server-client.key"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "etcd", "client.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "etcd", "client.key"))
		case cloudControllerComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-"+version.Program+"-cloud-controller.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-"+version.Program+"-cloud-controller.key"))
		case version.Program + programControllerComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-"+version.Program+"-controller.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-"+version.Program+"-controller.key"))
			// agent dir
			certList = append(certList, filepath.Join(agentDataDir, "client-"+version.Program+"-controller.crt"))
			certList = append(certList, filepath.Join(agentDataDir, "client-"+version.Program+"-controller.key"))
		case authProxyComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-auth-proxy.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-auth-proxy.key"))
		case kubeletComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-kubelet.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-kubelet.key"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "serving-kubelet.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "serving-kubelet.key"))
			// agent dir
			certList = append(certList, filepath.Join(agentDataDir, "client-kubelet.crt"))
			certList = append(certList, filepath.Join(agentDataDir, "client-kubelet.key"))
			certList = append(certList, filepath.Join(agentDataDir, "serving-kubelet.crt"))
			certList = append(certList, filepath.Join(agentDataDir, "serving-kubelet.key"))
		case kubeProxyComponent:
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-kube-proxy.crt"))
			certList = append(certList, filepath.Join(serverConfig.ControlConfig.DataDir, "tls", "client-kube-proxy.key"))
			// agent dir
			certList = append(certList, filepath.Join(agentDataDir, "client-kube-proxy.crt"))
			certList = append(certList, filepath.Join(agentDataDir, "client-kube-proxy.key"))
		}
	}

	for _, cert := range certList {
		if err := os.Remove(cert); err == nil {
			logrus.Infof("Certificate %s is deleted", cert)
		}
	}
	return nil
}

func rotateAllCerts(dirs ...string) error {
	for _, dir := range dirs {
		err := filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if (strings.HasSuffix(path, ".crt") || strings.HasSuffix(path, "key")) &&
					!strings.Contains(path, "-ca") &&
					!strings.Contains(path, "service.key") &&
					!strings.Contains(path, "temporary-certs") {
					if err := os.Remove(path); err == nil {
						logrus.Infof("Certificate %s is deleted", path)
					}
					return nil
				}
				return nil
			})
		if err != nil {
			return err
		}
	}
	return nil
}
