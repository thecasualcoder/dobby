package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/config"
	"github.com/thecasualcoder/dobby/pkg/utils"
)

// Handler is provides HandlerFunc for Gin context
type Handler struct {
	isHealthy bool
	isReady   bool
}

// New creates a new Handler
func New(initialHealth, initialReadiness bool) *Handler {
	return &Handler{
		isReady:   initialHealth,
		isHealthy: initialReadiness,
	}
}

// Health return the dobby health status
func (h *Handler) Health(c *gin.Context) {
	statusCode := http.StatusOK
	if !h.isHealthy {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"healthy": h.isHealthy})
}

// Ready return the dobby health status
func (h *Handler) Ready(c *gin.Context) {
	statusCode := http.StatusOK
	if !h.isReady {
		statusCode = http.StatusServiceUnavailable
	}
	c.JSON(statusCode, gin.H{"ready": h.isReady})
}

// Version return dobby version
func (h *Handler) Version(c *gin.Context) {
	if !h.isReady {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "application is not ready"})
		return
	}
	if !h.isHealthy {
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

// Meta return dobby's metadata
func (h *Handler) Meta(c *gin.Context) {
	if !h.isReady {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "application is not ready"})
		return
	}
	if !h.isHealthy {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "application is not healthy"})
		return
	}
	ip, err := utils.GetOutboundIP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"IP": ip, "HostName": os.Getenv("HOSTNAME")})
}

// MakeHealthPerfect will make dobby's health perfect
func (h *Handler) MakeHealthPerfect(c *gin.Context) {
	h.isHealthy = true
	c.JSON(200, gin.H{"status": "success"})
}

// MakeHealthSick will make dobby's health sick
func (h *Handler) MakeHealthSick(c *gin.Context) {
	h.isHealthy = false
	c.JSON(200, gin.H{"status": "success"})
}

// MakeReadyPerfect will make dobby's readiness perfect
func (h *Handler) MakeReadyPerfect(c *gin.Context) {
	h.isReady = true
	c.JSON(200, gin.H{"status": "success"})
}

// MakeReadySick will make dobby's readiness sick
func (h *Handler) MakeReadySick(c *gin.Context) {
	h.isReady = false
	c.JSON(200, gin.H{"status": "success"})
}

func init() {
}
