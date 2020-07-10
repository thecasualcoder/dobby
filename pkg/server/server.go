package server

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/handler"
)

// Run the gin server in address and port specified
func Run(bindAddress, port string, initialHealth, initialReadiness bool) error {
	r := gin.Default()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", bindAddress, port),
		Handler: r,
	}

	Bind(r, server, initialHealth, initialReadiness)
	return server.ListenAndServe()
}

// Bind binds all the routes to gin engine
func Bind(root *gin.Engine, server *http.Server, initialHealth, initialReadiness bool) {
	h := handler.New(initialHealth, initialReadiness)
	{
		root.GET("/health", h.Health)
		root.GET("/readiness", h.Ready)
		root.GET("/version", h.Version)
		root.GET("/meta", h.Meta)
		root.GET("/return/:statusCode", h.HTTPStat)
		root.POST("/call", h.Call)
	}
	controlGroup := root.Group("/control")
	{
		controlGroup.PUT("/health/perfect", h.MakeHealthPerfect)
		controlGroup.PUT("/health/sick", h.MakeHealthSick)
		controlGroup.PUT("/ready/perfect", h.MakeReadyPerfect)
		controlGroup.PUT("/ready/sick", h.MakeReadySick)
		controlGroup.PUT("/goturbo/memory", handler.GoTurboMemory)
		controlGroup.PUT("/goturbo/cpu", handler.GoTurboCPU)
		controlGroup.PUT("/crash", handler.Crash(server))
	}
}
