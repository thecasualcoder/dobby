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
	return r.Run(fmt.Sprintf("%s:%s", bindAddress, port))
}
