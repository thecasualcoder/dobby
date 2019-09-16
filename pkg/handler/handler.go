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
	isReady   = true
)

// Health return the dobby health status
func Health(c *gin.Context) {
	statusCode := http.StatusOK
	if !isHealthy {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"healthy": isHealthy})
}

// Ready return the dobby health status
func Ready(c *gin.Context) {
	statusCode := http.StatusOK
	if !isReady {
		statusCode = http.StatusServiceUnavailable
	}
	c.JSON(statusCode, gin.H{"ready": isReady})
}

// Version return dobby version
func Version(c *gin.Context) {
	if !isReady {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "application is not ready"})
		return
	}
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

// MakeReadyPerfect will make dobby's readiness perfect
func MakeReadyPerfect(c *gin.Context) {
	isReady = true
	c.JSON(200, gin.H{"status": "success"})
}

// MakeReadySick will make dobby's readiness sick
func MakeReadySick(c *gin.Context) {
	isReady = false
	c.JSON(200, gin.H{"status": "success"})
}

func init() {
	if initialHealth, err := strconv.ParseBool(os.Getenv("INITIAL_HEALTH")); err == nil {
		isHealthy = initialHealth
	}

	if initialReadiness, err := strconv.ParseBool(os.Getenv("INITIAL_READINESS")); err == nil {
		isReady = initialReadiness
	}
}
