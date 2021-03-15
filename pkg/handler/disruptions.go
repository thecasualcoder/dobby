package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/model"
	"log"
	"net/http"
)

// Crash will make dobby to kill itself
// As dobby dies, the gin server also shuts down.
// @Summary Suicide
// @Description Make Dobby kill itself
// @Tags Control
// @Accept json
// @Router /control/crash [put]
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
// @Summary Memory Spike
// @Description Make Dobby create a memory spike
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} model.ControlSuccess
// @Router /control/goturbo/memory [put]
func GoTurboMemory(c *gin.Context) {
	memorySpike := []string{"qwertyuiopasdfghjklzxcvbnm"}
	go func() {
		for {
			memorySpike = append(memorySpike, memorySpike...)
		}
	}()
	c.JSON(200, model.ControlSuccess{Status: "success"})
}

// GoTurboCPU will make dobby go Turbo
// Watch the video `https://youtu.be/TNjAZZ3vQ8o?t=14`
// for more context on `Going Turbo`
// @Summary CPU Spike
// @Description Make Dobby create a CPU spike
// @Tags Control
// @Accept json
// @Produce json
// @Success 200 {object} model.ControlSuccess
// @Router /control/goturbo/cpu [put]
func GoTurboCPU(c *gin.Context) {
	go func() {
		for {
			_ = 0
		}
	}()
	c.JSON(200, model.ControlSuccess{Status: "success"})
}
