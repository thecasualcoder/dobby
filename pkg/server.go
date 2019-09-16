package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/handler"
	"log"
	"net/http"
)

// Run the gin server in address and port specified
func Run(bindAddress, port string) error {
	r := gin.Default()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", bindAddress, port),
		Handler: r,
	}

	{
		r.GET("/health", handler.Health)
		r.GET("/readiness", handler.Ready)
		r.GET("/version", handler.Version)
	}

	controlGroup := r.Group("/control")
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

	return server.ListenAndServe()
}
