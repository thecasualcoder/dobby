package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/model"
	"net/http"
	"strconv"
	"time"
)

// HTTPStat returns the status code send by the client
// @Summary Repeat Status
// @Description Ask Dobby to return the status code sent by the client
// @Tags Status
// @Accept json
// @Produce json
// @Failure 400 {object} model.Error
// @Param statusCode path int true "Status Code - E.g. 200"
// @Param delay query int false "Dela(milliseconds) - E.g. 1000"
// @Router /return/{statusCode} [get]
func (h *Handler) HTTPStat(c *gin.Context) {
	returnCodeStr := c.Param("statusCode")
	returnCode, err := strconv.Atoi(returnCodeStr)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.Error{Error: fmt.Sprintf("error converting the statusCode to int: %s", err.Error())},
		)
		return
	}

	if delayStr := c.Query("delay"); delayStr != "" {
		delay, err := strconv.Atoi(delayStr)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				model.Error{Error: fmt.Sprintf("error converting the delay to int: %s", err.Error())},
			)
			return
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	c.Status(returnCode)
}
