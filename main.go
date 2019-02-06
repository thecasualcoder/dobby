package main

import (
	"github.com/thecasualcoder/dobby/cmd"
	"log"
	"os"
)

var majorVersion string
var minorVersion string

func main() {
	cmd.SetBuildVersion(majorVersion, minorVersion)
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
