package main

import (
	"github.com/thecasualcoder/dobby/cmd"
	"github.com/thecasualcoder/dobby/pkg/config"
	"log"
	"os"
)

var majorVersion string
var minorVersion string

func main() {
	config.SetBuildVersion(majorVersion, minorVersion)
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
