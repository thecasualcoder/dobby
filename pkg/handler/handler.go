package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler is provides HandlerFunc for Gin Context
type Handler struct {
	isHealthy     bool
	isReady       bool
	client        httpClient
	proxyRequests proxyRequests
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// New creates a new Handler
func New(initialHealth, initialReadiness bool, httpClient httpClient) *Handler {
	return &Handler{
		isReady:       initialReadiness,
		isHealthy:     initialHealth,
		client:        httpClient,
		proxyRequests: make(proxyRequests, 0),
	}
}

// Context is the interface represents the minimalistic gin.Context
// this is used to create mock struct while testing
type Context interface {
	JSON(code int, obj interface{})
	GetRequestBody() io.ReadCloser
	Status(code int)
	GetURI() *url.URL
	GetMethod() string
	SendResponse(response *http.Response, url string)
}

// NewDefaultContext creates the wrapper Context with gin Context
func NewDefaultContext(c *gin.Context) Context {
	return defaultContext{ginContext: c}
}

type defaultContext struct {
	ginContext *gin.Context
}

func (c defaultContext) GetURI() *url.URL {
	return c.ginContext.Request.URL
}

func (c defaultContext) GetMethod() string {
	return c.ginContext.Request.Method
}

func (c defaultContext) Status(code int) {
	c.ginContext.Status(code)
}

func (c defaultContext) JSON(code int, obj interface{}) {
	c.ginContext.JSON(code, obj)
}

func (c defaultContext) GetRequestBody() io.ReadCloser {
	return c.ginContext.Request.Body
}

func (c defaultContext) SendResponse(response *http.Response, url string) {
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when reading response from %s: %s", url, err.Error())})
		return
	}
	if len(responseData) == 0 {
		c.Status(response.StatusCode)
		return
	}
	responseStr := string(responseData)
	if strings.HasPrefix(responseStr, "{") || strings.HasPrefix(responseStr, "[") {
		var res interface{}
		err = json.Unmarshal(responseData, &res)
		if err != nil {
			c.JSON(400, gin.H{"error": fmt.Sprintf("error when decoding response from %s: %s", url, err.Error())})
			return
		}
		c.JSON(response.StatusCode, res)
		return
	}
	c.JSON(response.StatusCode, responseStr)
}
