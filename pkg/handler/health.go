package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/model"
)

// Health return the dobby health status
// @Summary Dobby Health
// @Description Get Dobby's health status
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} model.Health
// @Failure 500 {object} model.Health
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	statusCode := http.StatusOK
	if !h.isHealthy {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, model.Health{Healthy: h.isHealthy})
}

// MakeHealthPerfect godoc
// @Summary Make Healthy
// @Description Make Dobby healthy
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} model.ControlSuccess
// @Router /control/health/perfect [put]
func (h *Handler) MakeHealthPerfect(c *gin.Context) {
	h.isHealthy = true
	c.JSON(200, model.ControlSuccess{Status: "success"})
}

// MakeHealthSick godoc
// @Summary Make Unhealthy
// @Description Make Dobby sick or unhealthy
// @Tags Control
// @Accept json
// @Produce json
// @Param resetInSeconds query int false "Recover health after sometime (seconds) - E.g. 2"
// @Success 200 {object} model.ControlSuccess
// @Router /control/health/sick [put]
func (h *Handler) MakeHealthSick(c *gin.Context) {
	h.isHealthy = false
	setupResetFunction(c, func() {
		h.isHealthy = true
	})
	c.JSON(200, model.ControlSuccess{Status: "success"})
}

func setupResetFunction(c *gin.Context, afterFunc func()) {
	const resetInSecondsQueryParam = "resetInSeconds"
	resetTimer := c.Query(resetInSecondsQueryParam)
	if resetInSeconds, err := strconv.Atoi(resetTimer); err == nil && resetInSeconds != 0 {
		go time.AfterFunc(time.Second*time.Duration(resetInSeconds), afterFunc)
	}
}
