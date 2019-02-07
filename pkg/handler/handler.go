package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/config"
	"os"
)

// Health return the dobby health status
func Health(c *gin.Context) {
	c.JSON(200, gin.H{"healthy": true})
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
