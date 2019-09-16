package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/config"
	"net/http"
	"os"
	"strconv"
)

var (
	isHealthy = true
)

// Health return the dobby health status
func Health(c *gin.Context) {
	statusCode := http.StatusOK
	if !isHealthy {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"healthy": isHealthy})
}

// Version return dobby version
func Version(c *gin.Context) {
	if !isHealthy {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "application is not healthy"})
		return
	}

	envVersion := os.Getenv("VERSION")
	version := config.BuildVersion()
	if envVersion != "" {
		version = envVersion
	}
	c.JSON(200, gin.H{"version": version})
}

// MakeHealthPerfect will make dobby's health perfect
func MakeHealthPerfect(c *gin.Context) {
	isHealthy = true
	c.JSON(200, gin.H{"status": "success"})
}

// MakeHealthSick will make dobby's health sick
func MakeHealthSick(c *gin.Context) {
	isHealthy = false
	c.JSON(200, gin.H{"status": "success"})
}

func init() {
	if initialHealth, err := strconv.ParseBool(os.Getenv("INITIAL_HEALTH")); err == nil {
		isHealthy = initialHealth
	}
}
