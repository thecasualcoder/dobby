package cmd

import (
	"github.com/thecasualcoder/dobby/pkg/server"
	"github.com/urfave/cli"
)

func serverCmd() cli.Command {
	return cli.Command{
		Name:        "server",
		Description: "run dobby in server mode",
		Usage:       "run dobby server",
		Action:      runServer,
		Flags:       serverFlags(),
	}
}

func serverFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "bind-address, a",
			Value:  "127.0.0.1",
			Usage:  "Address of dobby server.",
			EnvVar: "BIND_ADDR",
		},
		cli.StringFlag{
			Name:   "port, p",
			Value:  "4444",
			Usage:  "Port which will be used by dobby server.",
			EnvVar: "PORT",
		},
	}
}

func runServer(context *cli.Context) {
	bindAddress := context.String("bind-address")
	port := context.String("port")
	err := server.Run(bindAddress, port)
	dieIf(err)
}
