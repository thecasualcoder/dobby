package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/config"
	"github.com/thecasualcoder/dobby/pkg/utils"
	"strconv"
	"time"
)

// Handler is provides HandlerFunc for Gin context
type Handler struct {
	isHealthy bool
	isReady   bool
}

// New creates a new Handler
func New(initialHealth, initialReadiness bool) *Handler {
	return &Handler{
		isReady:   initialReadiness,
		isHealthy: initialHealth,
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

// HTTPStat returns the status code send by the client
func (h *Handler) HTTPStat(c *gin.Context) {
	returnCodeStr := c.Param("statusCode")
	returnCode, err := strconv.Atoi(returnCodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error converting the statusCode to int: %s", err.Error())})
		return
	}

	if delayStr := c.Query("delay"); delayStr != "" {
		delay, err := strconv.Atoi(delayStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error converting the delay to int: %s", err.Error())})
			return
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	c.Status(returnCode)
}

// MakeHealthPerfect will make dobby's health perfect
func (h *Handler) MakeHealthPerfect(c *gin.Context) {
	h.isHealthy = true
	c.JSON(200, gin.H{"status": "success"})
}

// MakeHealthSick will make dobby's health sick
func (h *Handler) MakeHealthSick(c *gin.Context) {
	h.isHealthy = false
	setupResetFunction(c, func() {
		h.isHealthy = true
	})
	c.JSON(200, gin.H{"status": "success"})
}

func setupResetFunction(c *gin.Context, afterFunc func()) {
	const resetInSecondsQueryParam = "resetInSeconds"
	resetTimer := c.Query(resetInSecondsQueryParam)
	if resetInSeconds, err := strconv.Atoi(resetTimer); err == nil && resetInSeconds != 0 {
		go time.AfterFunc(time.Second*time.Duration(resetInSeconds), afterFunc)
	}
}

// MakeReadyPerfect will make dobby's readiness perfect
func (h *Handler) MakeReadyPerfect(c *gin.Context) {
	h.isReady = true
	c.JSON(200, gin.H{"status": "success"})
}

// MakeReadySick will make dobby's readiness sick
func (h *Handler) MakeReadySick(c *gin.Context) {
	h.isReady = false
	setupResetFunction(c, func() {
		h.isReady = true
	})
	c.JSON(200, gin.H{"status": "success"})
}

// Crash will make dobby to kill itself
// As dobby dies, the gin server also shuts down.
func Crash(server *http.Server) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		defer func() {
			_ = server.Shutdown(ctx)
		}()

		log.Fatal("you asked me do so, killing myself :-)")
	}
}

// GoTurboMemory will make dobby go Turbo
// Watch the video `https://youtu.be/TNjAZZ3vQ8o?t=14`
// for more context on `Going Turbo`
func GoTurboMemory(c *gin.Context) {
	memorySpike := []string{"qwertyuiopasdfghjklzxcvbnm"}
	go func() {
		for {
			memorySpike = append(memorySpike, memorySpike...)
		}
	}()
	c.JSON(200, gin.H{"status": "success"})
}

// GoTurboCPU will make dobby go Turbo
// Watch the video `https://youtu.be/TNjAZZ3vQ8o?t=14`
// for more context on `Going Turbo`
func GoTurboCPU(c *gin.Context) {
	go func() {
		for {
			_ = 0
		}
	}()
	c.JSON(200, gin.H{"status": "success"})
}
