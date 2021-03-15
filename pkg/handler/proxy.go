package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
