package cmds

import (
	"github.com/rancher/k3s/pkg/version"
	"github.com/urfave/cli"
)

const CertCommand = "cert"

var CertCommandFlags = []cli.Flag{
	DebugFlag,
	ConfigFlag,
	LogFile,
	AlsoLogToStderr,
	cli.StringFlag{
		Name:        "data-dir,d",
		Usage:       "(data) Folder to hold state default /var/lib/rancher/" + version.Program + " or ${HOME}/.rancher/" + version.Program + " if not root",
		Destination: &ServerConfig.DataDir,
	},
}

func NewCertCommand(subcommands []cli.Command) cli.Command {
	return cli.Command{
		Name:            CertCommand,
		Usage:           "Certificates management",
		SkipFlagParsing: false,
		SkipArgReorder:  true,
		Subcommands:     subcommands,
		Flags:           CertCommandFlags,
	}
}

func NewCertSubcommands(rotate func(ctx *cli.Context) error) []cli.Command {
	return []cli.Command{
		{
			Name:            "rotate",
			Usage:           "Certificate rotation",
			SkipFlagParsing: false,
			SkipArgReorder:  true,
			Action:          rotate,
			Flags:           CertCommandFlags,
		},
	}
}
