package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/model"
	"github.com/thecasualcoder/dobby/pkg/utils"
	"net/http"
	"os"
)

// Meta return dobby's metadata
// @Summary Dobby Metadata
// @Description Get Dobby's metadata
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} model.Metadata
// @Failure 503 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /meta [get]
func (h *Handler) Meta(c *gin.Context) {
	if !h.isReady {
		c.JSON(http.StatusServiceUnavailable, model.Error{Error: "application is not ready"})
		return
	}
	if !h.isHealthy {
		c.JSON(http.StatusInternalServerError, model.Error{Error: "application is not healthy"})
		return
	}
	ip, err := utils.GetOutboundIP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Metadata{IP: ip, Hostname: os.Getenv("HOSTNAME")})
}
