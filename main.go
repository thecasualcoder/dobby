package main

import (
	"github.com/thecasualcoder/dobby/cmd"
	"github.com/thecasualcoder/dobby/docs"
	"github.com/thecasualcoder/dobby/pkg/config"
	"log"
	"os"
)

var majorVersion string
var minorVersion string

func main() {
	config.SetBuildVersion(majorVersion, minorVersion)

	docs.SwaggerInfo.Title = "Dobby"
	docs.SwaggerInfo.Description = "dobby is free and will serve your orders."
	docs.SwaggerInfo.Version = config.BuildVersion()
	docs.SwaggerInfo.Schemes = []string{"http"}

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
