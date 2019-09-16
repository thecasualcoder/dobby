package server

import (
	"fmt"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/handler"
)

// Run the gin server in address and port specified
func Run(bindAddress, port string) error {
	r := gin.Default()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", bindAddress, port),
		Handler: r,
	}

	Bind(r, server)
	return server.ListenAndServe()
}

// Bind binds all the routes to gin engine
func Bind(root *gin.Engine, server *http.Server) {
	{
		root.GET("/health", handler.Health)
		root.GET("/readiness", handler.Ready)
		root.GET("/meta", handler.Meta)
	}
	controlGroup := root.Group("/control")
	{
		controlGroup.PUT("/health/perfect", handler.MakeHealthPerfect)
		controlGroup.PUT("/health/sick", handler.MakeHealthSick)
		controlGroup.PUT("/ready/perfect", handler.MakeReadyPerfect)
		controlGroup.PUT("/ready/sick", handler.MakeReadySick)
		controlGroup.PUT("/crash", func(ctx *gin.Context) {
			defer func() {
				_ = server.Shutdown(ctx)
			}()

			log.Fatal("you asked me do so, killing myself :-)")
		})
	}
}
