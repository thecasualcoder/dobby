package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thecasualcoder/dobby/pkg/config"
	"github.com/thecasualcoder/dobby/pkg/model"
	"github.com/thecasualcoder/dobby/pkg/utils"
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
// @Summary Dobby Health
// @Description Get Dobby's health status
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} model.Health
// @Failure 500 {object} model.Health
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	statusCode := http.StatusOK
	if !h.isHealthy {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, model.Health{Healthy: h.isHealthy})
}

// Ready return the dobby health status
// @Summary Dobby Ready
// @Description Get Dobby's readiness
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} model.Ready
// @Failure 503 {object} model.Ready
// @Router /ready [get]
func (h *Handler) Ready(c *gin.Context) {
	statusCode := http.StatusOK
	if !h.isReady {
		statusCode = http.StatusServiceUnavailable
	}
	c.JSON(statusCode, model.Ready{Ready: h.isReady})
}

// Version return dobby version
// @Summary Dobby Version
// @Description Get Dobby's version
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {object} model.Version
// @Failure 503 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /version [get]
func (h *Handler) Version(c *gin.Context) {
	if !h.isReady {
		c.JSON(http.StatusServiceUnavailable, model.Error{Error: "application is not ready"})
		return
	}
	if !h.isHealthy {
		c.JSON(http.StatusInternalServerError, model.Error{Error: "application is not healthy"})
		return
	}

	envVersion := os.Getenv("VERSION")
	version := config.BuildVersion()
	if envVersion != "" {
		version = envVersion
	}
	c.JSON(200, model.Version{Version: version})
}

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

// ProxyRoute will route to custom route if the route is found in proxyRequests
// this will be invoked when no standard routes are found in gin
func (h *Handler) ProxyRoute(c Context) {
	proxyConfig := h.proxyRequests.getProxy(c.GetURI().Path, c.GetMethod())
	if proxyConfig == nil {
		c.Status(404)
		return
	}
	request, err := http.NewRequest(proxyConfig.Method, proxyConfig.URL, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when creating request for %s: %v", proxyConfig.URL, err.Error())})
		return
	}
	response, err := h.client.Do(request)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when creating request for %s: %s", proxyConfig.URL, proxyConfig.Method)})
		return
	}
	sendResponse(c, response, proxyConfig.URL)
}

func sendResponse(c Context, response *http.Response, url string) {
	responseData, err := ioutil.ReadAll(response.Body)
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

// AddProxy will add the proxy settings
func (h *Handler) AddProxy(c Context) {
	decoder := json.NewDecoder(c.GetRequestBody())
	var proxyRequest proxyRequest
	err := decoder.Decode(&proxyRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when decoding request: %s", err.Error())})
		return
	}
	if h.proxyRequests.isPresent(proxyRequest) {
		c.JSON(400, gin.H{"error": fmt.Sprintf("proxy configuration for url: %s and method: %s is already added", proxyRequest.Path, proxyRequest.Method)})
		return
	}
	h.proxyRequests = append(h.proxyRequests, proxyRequest)
	c.Status(201)
}

// DeleteProxy will delete the proxy configuration
func (h *Handler) DeleteProxy(c Context) {
	decoder := json.NewDecoder(c.GetRequestBody())
	var proxyRequest proxyRequest
	err := decoder.Decode(&proxyRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("error when decoding request: %s", err.Error())})
		return
	}
	if h.proxyRequests.isPresent(proxyRequest) {
		h.proxyRequests = h.proxyRequests.deleteProxy(proxyRequest.Path, proxyRequest.Method)
		c.JSON(200, gin.H{"result": "deleted the proxy config successfully"})
		return
	}
	c.JSON(404, gin.H{"error": fmt.Sprintf("proxy config with url %s and %s method is not found", proxyRequest.Path, proxyRequest.Method)})
}

type proxyRequest struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Proxy  proxy  `json:"proxy"`
}

type proxyRequests []proxyRequest

func (ps proxyRequests) isPresent(requestedProxyRequest proxyRequest) bool {
	for _, p := range ps {
		if (p.Path == requestedProxyRequest.Path) && (p.Method == requestedProxyRequest.Method) {
			return true
		}
	}
	return false
}

func (ps proxyRequests) getProxy(path string, method string) *proxy {
	for _, p := range ps {
		if p.Path == path && p.Method == method {
			return &p.Proxy
		}
	}
	return nil
}

func (ps proxyRequests) deleteProxy(path string, method string) proxyRequests {
	accProxyRequests := make(proxyRequests, 0, len(ps))
	for _, p := range ps {
		if p.Path != path || p.Method != method {
			accProxyRequests = append(accProxyRequests, p)
		}
	}
	return accProxyRequests
}

type proxy struct {
	URL    string `json:"url"`
	Method string `json:"method"`
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
