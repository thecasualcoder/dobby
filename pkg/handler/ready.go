package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/model"
	"net/http"
)

// Ready return the dobby health status
// @Summary Dobby Ready
// @Description Get Dobby's readiness
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} model.Ready
// @Failure 503 {object} model.Ready
// @Router /ready [get]
func (h *Handler) Ready(c *gin.Context) {
	statusCode := http.StatusOK
	if !h.isReady {
		statusCode = http.StatusServiceUnavailable
	}
	c.JSON(statusCode, model.Ready{Ready: h.isReady})
}

// MakeReadyPerfect godoc
// @Summary Make Ready
// @Description Make Dobby ready
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} model.ControlSuccess
// @Router /control/ready/perfect [put]
func (h *Handler) MakeReadyPerfect(c *gin.Context) {
	h.isReady = true
	c.JSON(200, model.ControlSuccess{Status: "success"})
}

// MakeReadySick godoc
// @Summary Make Unready
// @Description Make Dobby unready
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} model.ControlSuccess
// @Param resetInSeconds query int false "Recover readiness after sometime (seconds) - E.g. 2"
// @Router /control/ready/sick [put]
func (h *Handler) MakeReadySick(c *gin.Context) {
	h.isReady = false
	setupResetFunction(c, func() {
		h.isReady = true
	})
	c.JSON(200, model.ControlSuccess{Status: "success"})
}
