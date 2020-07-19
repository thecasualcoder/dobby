package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/config"
	"github.com/thecasualcoder/dobby/pkg/utils"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Handler is provides HandlerFunc for Gin Context
type Handler struct {
	isHealthy     bool
	isReady       bool
	client        httpClient
	proxyRequests ProxyRequests
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
		proxyRequests: make(ProxyRequests, 0),
	}
}

// Context is the interface represents the minimalistic gin.Context
// this is used to create mock struct while testing
type Context interface {
	Header(key, value string)
	Data(code int, contentType string, data []byte)
	JSON(code int, obj interface{})
	GetRequestBody() io.ReadCloser
	GetRequestHeader() http.Header
	Status(code int)
	GetURI() *url.URL
	GetMethod() string
}

type defaultContext struct {
	ginContext *gin.Context
}

func (c defaultContext) GetRequestHeader() http.Header {
	return c.ginContext.Request.Header
}

func (c defaultContext) Header(key, value string) {
	c.ginContext.Header(key, value)
}

func (c defaultContext) Data(code int, contentType string, data []byte) {
	c.ginContext.Data(code, contentType, data)
}

func (c defaultContext) GetURI() *url.URL {
	return c.ginContext.Request.URL
}

func (c defaultContext) GetMethod() string {
	return c.ginContext.Request.Method
}

// NewDefaultContext creates the wrapper Context with gin Context
func NewDefaultContext(c *gin.Context) Context {
	return defaultContext{ginContext: c}
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

// Health return the dobby health status
func (h *Handler) Health(c *gin.Context) {
	statusCode := http.StatusOK
	if !h.isHealthy {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"healthy": h.isHealthy})
}

// Ready return the dobby health status
func (h *Handler) Ready(c *gin.Context) {
	statusCode := http.StatusOK
	if !h.isReady {
		statusCode = http.StatusServiceUnavailable
	}
	c.JSON(statusCode, gin.H{"ready": h.isReady})
}

// Version return dobby version
func (h *Handler) Version(c *gin.Context) {
	if !h.isReady {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "application is not ready"})
		return
	}
	if !h.isHealthy {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "application is not healthy"})
		return
	}

	envVersion := os.Getenv("VERSION")
	version := config.BuildVersion()
	if envVersion != "" {
		version = envVersion
	}
	c.JSON(200, gin.H{"version": version})
}

// Meta return dobby's metadata
func (h *Handler) Meta(c *gin.Context) {
	if !h.isReady {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "application is not ready"})
		return
	}
	if !h.isHealthy {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "application is not healthy"})
		return
	}
	ip, err := utils.GetOutboundIP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"IP": ip, "HostName": os.Getenv("HOSTNAME")})
}

// HTTPStat returns the status code send by the client
func (h *Handler) HTTPStat(c *gin.Context) {
	returnCodeStr := c.Param("statusCode")
	returnCode, err := strconv.Atoi(returnCodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error converting the statusCode to int: %s", err.Error())})
		return
	}

	if delayStr := c.Query("delay"); delayStr != "" {
		delay, err := strconv.Atoi(delayStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error converting the delay to int: %s", err.Error())})
			return
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	c.Status(returnCode)
}

// MakeHealthPerfect will make dobby's health perfect
func (h *Handler) MakeHealthPerfect(c *gin.Context) {
	h.isHealthy = true
	c.JSON(200, gin.H{"status": "success"})
}

// MakeHealthSick will make dobby's health sick
func (h *Handler) MakeHealthSick(c *gin.Context) {
	h.isHealthy = false
	setupResetFunction(c, func() {
		h.isHealthy = true
	})
	c.JSON(200, gin.H{"status": "success"})
}

func setupResetFunction(c *gin.Context, afterFunc func()) {
	const resetInSecondsQueryParam = "resetInSeconds"
	resetTimer := c.Query(resetInSecondsQueryParam)
	if resetInSeconds, err := strconv.Atoi(resetTimer); err == nil && resetInSeconds != 0 {
		go time.AfterFunc(time.Second*time.Duration(resetInSeconds), afterFunc)
	}
}

// MakeReadyPerfect will make dobby's readiness perfect
func (h *Handler) MakeReadyPerfect(c *gin.Context) {
	h.isReady = true
	c.JSON(200, gin.H{"status": "success"})
}

// MakeReadySick will make dobby's readiness sick
func (h *Handler) MakeReadySick(c *gin.Context) {
	h.isReady = false
	setupResetFunction(c, func() {
		h.isReady = true
	})
	c.JSON(200, gin.H{"status": "success"})
}

// Call another service and send the response
func (h *Handler) Call(c Context) {
	decoder := json.NewDecoder(c.GetRequestBody())
	var callRequest callRequest
	err := decoder.Decode(&callRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when decoding request: %s", err.Error())})
		return
	}
	response, err := h.makeCall(callRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when making request to %s: %s", callRequest.URL, err.Error())})
		return
	}
	sendResponse(c, response, callRequest.URL)
}

func (h *Handler) makeCall(callRequest callRequest) (*http.Response, error) {
	marshal, err := json.Marshal(callRequest.Body)
	if err != nil {
		return nil, fmt.Errorf("error when marshalling request body: %s", err)
	}
	request, err := http.NewRequest(callRequest.Method, callRequest.URL, bytes.NewBuffer(marshal))
	if err != nil {
		return nil, fmt.Errorf("error when creating new request to %s: %s", callRequest.URL, err)
	}
	return h.client.Do(request)
}

// ProxyRoute will route to custom route if the route is found in ProxyRequests
// this will be invoked when no standard routes are found in gin
func (h *Handler) ProxyRoute(c Context) {
	requestPath := c.GetURI().Path
	requestMethod := c.GetMethod()

	proxyConfig := h.proxyRequests.getProxy(requestPath, requestMethod)
	if proxyConfig == nil {
		proxyConfig = GetDynamicProxy(requestPath, requestMethod)
	}

	if proxyConfig == nil {
		c.Status(404)
		return
	}
	request, err := http.NewRequest(proxyConfig.Method, proxyConfig.URL, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when creating request for %s: %v", proxyConfig.URL, err.Error())})
		return
	}

	for key, values := range c.GetRequestHeader() {
		request.Header.Set(key, strings.Trim(strings.Join(values, "; "), "; "))
	}
	response, err := h.client.Do(request)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when creating request for %s: %s", proxyConfig.URL, proxyConfig.Method)})
		return
	}
	sendResponse(c, response, proxyConfig.URL)
}

func sendResponse(c Context, response *http.Response, url string) {
	for key, values := range response.Header {
		c.Header(key, strings.Trim(strings.Join(values, "; "), "; "))
	}
	w := bytes.Buffer{}
	_, _ = io.Copy(&w, response.Body)
	c.Data(response.StatusCode, response.Header.Get("Content-Type"), w.Bytes())
}

// AddProxy will add the Proxy settings
func (h *Handler) AddProxy(c Context) {
	decoder := json.NewDecoder(c.GetRequestBody())
	var proxyRequest ProxyRequest
	err := decoder.Decode(&proxyRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when decoding request: %s", err.Error())})
		return
	}
	if h.proxyRequests.isPresent(proxyRequest) {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Proxy configuration for url: %s and method: %s is already added", proxyRequest.Path, proxyRequest.Method)})
		return
	}
	h.proxyRequests = append(h.proxyRequests, proxyRequest)
	c.Status(201)
}

// DeleteProxy will delete the Proxy configuration
func (h *Handler) DeleteProxy(c Context) {
	decoder := json.NewDecoder(c.GetRequestBody())
	var proxyRequest ProxyRequest
	err := decoder.Decode(&proxyRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when decoding request: %s", err.Error())})
		return
	}
	if h.proxyRequests.isPresent(proxyRequest) {
		h.proxyRequests = h.proxyRequests.deleteProxy(proxyRequest.Path, proxyRequest.Method)
		c.JSON(200, gin.H{"result": "deleted the Proxy config successfully"})
		return
	}
	c.JSON(404, gin.H{"error": fmt.Sprintf("Proxy config with url %s and %s method is not found", proxyRequest.Path, proxyRequest.Method)})
}

// ProxyRequest represent a Proxy request
type ProxyRequest struct {
	Path   string `json:"path" yaml:"path"`
	Method string `json:"method" yaml:"method"`
	Proxy  Proxy  `json:"Proxy" yaml:"proxy"`
}

// ProxyRequests represents collection of Proxy requests
type ProxyRequests []ProxyRequest

// NewProxyRequests represents
func NewProxyRequests(data []byte) (ProxyRequests, error) {
	var result = ProxyRequests{}
	err := yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ps ProxyRequests) isPresent(requestedProxyRequest ProxyRequest) bool {
	for _, p := range ps {
		if (p.Path == requestedProxyRequest.Path) && (p.Method == requestedProxyRequest.Method) {
			return true
		}
	}
	return false
}

func (ps ProxyRequests) getProxy(path string, method string) *Proxy {
	for _, p := range ps {
		if p.Path == path && p.Method == method {
			return &p.Proxy
		}
	}
	return nil
}

func (ps ProxyRequests) deleteProxy(path string, method string) ProxyRequests {
	accProxyRequests := make(ProxyRequests, 0, len(ps))
	for _, p := range ps {
		if p.Path != path || p.Method != method {
			accProxyRequests = append(accProxyRequests, p)
		}
	}
	return accProxyRequests
}

// Proxy represents a proxy request
type Proxy struct {
	URL    string `json:"url" yaml:"url"`
	Method string `json:"method" yaml:"method"`
}

type callRequest struct {
	URL    string      `json:"url"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}

// Crash will make dobby to kill itself
// As dobby dies, the gin server also shuts down.
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
func GoTurboMemory(c *gin.Context) {
	memorySpike := []string{"qwertyuiopasdfghjklzxcvbnm"}
	go func() {
		for {
			memorySpike = append(memorySpike, memorySpike...)
		}
	}()
	c.JSON(200, gin.H{"status": "success"})
}

// GoTurboCPU will make dobby go Turbo
// Watch the video `https://youtu.be/TNjAZZ3vQ8o?t=14`
// for more context on `Going Turbo`
func GoTurboCPU(c *gin.Context) {
	go func() {
		for {
			_ = 0
		}
	}()
	c.JSON(200, gin.H{"status": "success"})
}
