package cmd

import (
	"github.com/thecasualcoder/dobby/pkg/config"
	"github.com/urfave/cli"
	"log"
)

func dieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Run is the entry point for dobby app
func Run(args []string) error {
	app := cli.NewApp()
	app.Description = "Web app which obey's invokers action"
	app.Usage = "Waiting to serve the order"
	app.Name = "dobby"
	app.Version = config.BuildVersion()

	app.Commands = []cli.Command{
		serverCmd(),
	}

	return app.Run(args)
}
