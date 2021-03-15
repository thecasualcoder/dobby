package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/config"
	"github.com/thecasualcoder/dobby/pkg/model"
	"net/http"
	"os"
)

// Version return dobby version
// @Summary Dobby Version
// @Description Get Dobby's version
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} model.Version
// @Failure 503 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /version [get]
func (h *Handler) Version(c *gin.Context) {
	if !h.isReady {
		c.JSON(http.StatusServiceUnavailable, model.Error{Error: "application is not ready"})
		return
	}
	if !h.isHealthy {
		c.JSON(http.StatusInternalServerError, model.Error{Error: "application is not healthy"})
		return
	}

	envVersion := os.Getenv("VERSION")
	version := config.BuildVersion()
	if envVersion != "" {
		version = envVersion
	}
	c.JSON(200, model.Version{Version: version})
}
