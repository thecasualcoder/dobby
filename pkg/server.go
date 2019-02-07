package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/handler"
)

// Run the gin server in address and port specified
func Run(bindAddress, port string) error {
	r := gin.Default()

	r.GET("/health", handler.Health)
	r.GET("/version", handler.Version)
	r.PUT("/state/crash", handler.Crash)
	return r.Run(fmt.Sprintf("%s:%s", bindAddress, port))
}
