// +build !no_etcd

package cmds

import (
	"github.com/urfave/cli"
)

func NewETCDCommand(action func(*cli.Context) error) cli.Command {
	return cli.Command{
		Name:      "etcd",
		Usage:     "Run etcd server",
		UsageText: appName + " etcd [OPTIONS]",
		Before:    SetupDebug(CheckSELinuxFlags),
		Action:    action,
		Flags: []cli.Flag{
			ConfigFlag,
			DebugFlag,
			VLevel,
			VModule,
			LogFile,
			AlsoLogToStderr,
			TLSSan,
			DataDir,
			Token,
			TokenFile,
			AgentToken,
			AgentTokenFile,
			KubeConfigOutput,
			KubeConfigMode,
			EtcdSnapshotDir,
			EtcdDisableSnapshots,
			EtcdSnapshotRetention,
			EtcdSnapshotCron,
			NodeNameFlag,
			WithNodeIDFlag,
			NodeLabels,
			NodeTaints,
			DockerFlag,
			CRIEndpointFlag,
			PauseImageFlag,
			SnapshotterFlag,
			PrivateRegistryFlag,
			NodeIPFlag,
			NodeExternalIPFlag,
			ResolvConfFlag,
			FlannelIfaceFlag,
			FlannelConfFlag,
			ExtraKubeletArgs,
			ExtraKubeProxyArgs,
			ProtectKernelDefaultsFlag,
			ServerURL,
			ClusterReset,
			ClusterResetRestorePath,
			&SELinuxFlag,

			// Hidden/Deprecated flags below

			&DisableSELinuxFlag,
			FlannelFlag,
		},
	}
}
