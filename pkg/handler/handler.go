package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/config"
	"log"
	"os"
	"strconv"
)

var statusCode = 200

// Health return the dobby health status
func Health(c *gin.Context) {
	healthy := true
	if statusCode != 200 {
		healthy = false
	}
	c.JSON(statusCode, gin.H{"healthy": healthy})
}

// Version return dobby version
func Version(c *gin.Context) {
	envVersion := os.Getenv("VERSION")
	version := config.BuildVersion()
	if envVersion != "" {
		version = envVersion
	}
	c.JSON(200, gin.H{"version": version})
}

// Crash will make dobby to kill itself
func Crash(_ *gin.Context) {
	log.Fatal("you asked me do so, killing myself :-)")
}

// Sick will make dobby sick
func Sick(c *gin.Context) {
	statusCode = 500
	c.JSON(200, gin.H{"status": "success"})
}

// Healthy will make dobby healthy again
func Healthy(c *gin.Context) {
	statusCode = 200
	c.JSON(200, gin.H{"status": "success"})
}

func init() {
	healthy, err := strconv.ParseBool(os.Getenv("HEALTH"))

	if err != nil {
		statusCode = 200
	} else if !healthy {
		statusCode = 500
	}
}
