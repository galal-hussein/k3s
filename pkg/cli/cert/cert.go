package cert

import (
	"io"
	"os"
	"path/filepath"

	"github.com/erikdubbelboer/gspt"
	"github.com/rancher/k3s/pkg/cli/cmds"
	"github.com/rancher/k3s/pkg/daemons/config"
	"github.com/rancher/k3s/pkg/server"
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

	// backing up tls dir
	if _, err := os.Stat(tlsDir); err != nil {
		return err
	}
	if err := io.Copy()
	// serverConfig.ControlConfig.Runtime.ClientCA = filepath.Join(dataDir, "tls", "client-ca.crt")
	// serverConfig.ControlConfig.Runtime.ClientCAKey = filepath.Join(dataDir, "tls", "client-ca.key")
	// serverConfig.ControlConfig.Runtime.ServerCA = filepath.Join(dataDir, "tls", "server-ca.crt")
	// serverConfig.ControlConfig.Runtime.ServerCAKey = filepath.Join(dataDir, "tls", "server-ca.key")
	// serverConfig.ControlConfig.Runtime.RequestHeaderCA = filepath.Join(dataDir, "tls", "request-header-ca.crt")
	// serverConfig.ControlConfig.Runtime.RequestHeaderCAKey = filepath.Join(dataDir, "tls", "request-header-ca.key")
	// serverConfig.ControlConfig.Runtime.IPSECKey = filepath.Join(dataDir, "cred", "ipsec.psk")

	// serverConfig.ControlConfig.Runtime.ServiceKey = filepath.Join(dataDir, "tls", "service.key")
	// serverConfig.ControlConfig.Runtime.PasswdFile = filepath.Join(dataDir, "cred", "passwd")
	// serverConfig.ControlConfig.Runtime.NodePasswdFile = filepath.Join(dataDir, "cred", "node-passwd")

	// serverConfig.ControlConfig.Runtime.KubeConfigAdmin = filepath.Join(dataDir, "cred", "admin.kubeconfig")
	// serverConfig.ControlConfig.Runtime.KubeConfigController = filepath.Join(dataDir, "cred", "controller.kubeconfig")
	// serverConfig.ControlConfig.Runtime.KubeConfigScheduler = filepath.Join(dataDir, "cred", "scheduler.kubeconfig")
	// serverConfig.ControlConfig.Runtime.KubeConfigAPIServer = filepath.Join(dataDir, "cred", "api-server.kubeconfig")
	// serverConfig.ControlConfig.Runtime.KubeConfigCloudController = filepath.Join(dataDir, "cred", "cloud-controller.kubeconfig")

	// serverConfig.ControlConfig.Runtime.ClientAdminCert = filepath.Join(dataDir, "tls", "client-admin.crt")
	// serverConfig.ControlConfig.Runtime.ClientAdminKey = filepath.Join(dataDir, "tls", "client-admin.key")
	// serverConfig.ControlConfig.Runtime.ClientControllerCert = filepath.Join(dataDir, "tls", "client-controller.crt")
	// serverConfig.ControlConfig.Runtime.ClientControllerKey = filepath.Join(dataDir, "tls", "client-controller.key")
	// serverConfig.ControlConfig.Runtime.ClientCloudControllerCert = filepath.Join(dataDir, "tls", "client-"+version.Program+"-cloud-controller.crt")
	// serverConfig.ControlConfig.Runtime.ClientCloudControllerKey = filepath.Join(dataDir, "tls", "client-"+version.Program+"-cloud-controller.key")
	// serverConfig.ControlConfig.Runtime.ClientSchedulerCert = filepath.Join(dataDir, "tls", "client-scheduler.crt")
	// serverConfig.ControlConfig.Runtime.ClientSchedulerKey = filepath.Join(dataDir, "tls", "client-scheduler.key")
	// serverConfig.ControlConfig.Runtime.ClientKubeAPICert = filepath.Join(dataDir, "tls", "client-kube-apiserver.crt")
	// serverConfig.ControlConfig.Runtime.ClientKubeAPIKey = filepath.Join(dataDir, "tls", "client-kube-apiserver.key")
	// serverConfig.ControlConfig.Runtime.ClientKubeProxyCert = filepath.Join(dataDir, "tls", "client-kube-proxy.crt")
	// serverConfig.ControlConfig.Runtime.ClientKubeProxyKey = filepath.Join(dataDir, "tls", "client-kube-proxy.key")
	// serverConfig.ControlConfig.Runtime.ClientK3sControllerCert = filepath.Join(dataDir, "tls", "client-"+version.Program+"-controller.crt")
	// serverConfig.ControlConfig.Runtime.ClientK3sControllerKey = filepath.Join(dataDir, "tls", "client-"+version.Program+"-controller.key")

	// serverConfig.ControlConfig.Runtime.ServingKubeAPICert = filepath.Join(dataDir, "tls", "serving-kube-apiserver.crt")
	// serverConfig.ControlConfig.Runtime.ServingKubeAPIKey = filepath.Join(dataDir, "tls", "serving-kube-apiserver.key")

	// serverConfig.ControlConfig.Runtime.ClientKubeletKey = filepath.Join(dataDir, "tls", "client-kubelet.key")
	// serverConfig.ControlConfig.Runtime.ServingKubeletKey = filepath.Join(dataDir, "tls", "serving-kubelet.key")

	// serverConfig.ControlConfig.Runtime.ClientAuthProxyCert = filepath.Join(dataDir, "tls", "client-auth-proxy.crt")
	// serverConfig.ControlConfig.Runtime.ClientAuthProxyKey = filepath.Join(dataDir, "tls", "client-auth-proxy.key")

	// serverConfig.ControlConfig.Runtime.ETCDServerCA = filepath.Join(dataDir, "tls", "etcd", "server-ca.crt")
	// serverConfig.ControlConfig.Runtime.ETCDServerCAKey = filepath.Join(dataDir, "tls", "etcd", "server-ca.key")
	// serverConfig.ControlConfig.Runtime.ETCDPeerCA = filepath.Join(dataDir, "tls", "etcd", "peer-ca.crt")
	// serverConfig.ControlConfig.Runtime.ETCDPeerCAKey = filepath.Join(dataDir, "tls", "etcd", "peer-ca.key")
	// serverConfig.ControlConfig.Runtime.ServerETCDCert = filepath.Join(dataDir, "tls", "etcd", "server-client.crt")
	// serverConfig.ControlConfig.Runtime.ServerETCDKey = filepath.Join(dataDir, "tls", "etcd", "server-client.key")
	// serverConfig.ControlConfig.Runtime.PeerServerClientETCDCert = filepath.Join(dataDir, "tls", "etcd", "peer-server-client.crt")
	// serverConfig.ControlConfig.Runtime.PeerServerClientETCDKey = filepath.Join(dataDir, "tls", "etcd", "peer-server-client.key")
	// serverConfig.ControlConfig.Runtime.ClientETCDCert = filepath.Join(dataDir, "tls", "etcd", "client.crt")
	// serverConfig.ControlConfig.Runtime.ClientETCDKey = filepath.Join(dataDir, "tls", "etcd", "client.key")

	// if config.EncryptSecrets {
	// 	serverConfig.ControlConfig.Runtime.EncryptionConfig = filepath.Join(dataDir, "cred", "encryption-config.json")
	// }
	// ctx := signals.SetupSignalHandler(context.Background())
	// // defining all certs that will not be rotated

	return nil
}
