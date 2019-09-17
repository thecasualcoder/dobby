package cmd

import (
	"github.com/thecasualcoder/dobby/pkg/server"
	"github.com/urfave/cli"
	"strconv"
	"time"
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
		cli.StringFlag{
			Name:   "initial-health",
			EnvVar: "INITIAL_HEALTH",
			Usage:  "Sets the Initial health of the server (/health) (true|false)",
			Value:  "true",
		},
		cli.StringFlag{
			Name:   "initial-readiness",
			EnvVar: "INITIAL_READINESS",
			Usage:  "Sets the Initial readiness of the server (/readiness) (true|false)",
			Value:  "true",
		},
		cli.Int64Flag{
			Name:   "initial-delay",
			EnvVar: "INITIAL_DELAY",
			Usage:  "Sets the Initial delay to start the server (in seconds)",
			Value:  0,
		},
	}
}

func runServer(context *cli.Context) {
	bindAddress := context.String("bind-address")
	port := context.String("port")
	initialDelay := time.Duration(context.Int64("initial-delay")) * time.Second
	time.Sleep(initialDelay)
	initialHealth := true
	if health, err := strconv.ParseBool(context.String("initial-health")); err == nil {
		initialHealth = health
	}

	initialReadiness := true
	if readiness, err := strconv.ParseBool(context.String("initial-readiness")); err == nil {
		initialReadiness = readiness
	}

	err := server.Run(bindAddress, port, initialHealth, initialReadiness)
	dieIf(err)
}
