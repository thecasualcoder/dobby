package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock "github.com/thecasualcoder/dobby/internal/mock/handler"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestHandler_AddProxy(t *testing.T) {
	t.Run("should add the proxy request to handler", func(t *testing.T) {
		handler := New(true, true, nil)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockContext := mock.NewMockContext(ctrl)
		stringReader := strings.NewReader(`
{
 "path": "/v1/version",
 "method": "GET",
 "proxy": {
   "url": "http://dobby2/version",
   "method": "GET"
 }
}`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		mockContext.EXPECT().Status(201)

		handler.AddProxy(mockContext)

		assert.Len(t, handler.proxyRequests, 1)
		expectedProxyRequest := proxyRequest{
			Path:   "/v1/version",
			Method: "GET", Proxy: proxy{
				URL:    "http://dobby2/version",
				Method: "GET",
			},
		}
		assert.Equal(t, expectedProxyRequest, handler.proxyRequests[0])
	})

	t.Run("should not add the proxy request if the same url and same method is added", func(t *testing.T) {
		handler := New(true, true, nil)
		handler.proxyRequests = proxyRequests{proxyRequest{
			Path:   "/v1/version",
			Method: "GET",
		}}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockContext := mock.NewMockContext(ctrl)
		stringReader := strings.NewReader(`
{
 "path": "/v1/version",
 "method": "GET"
}`)

		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		mockContext.EXPECT().JSON(400, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.EqualValues(t, "proxy configuration for url: /v1/version and method: GET is already added", data.(gin.H)["error"])
		})

		handler.AddProxy(mockContext)

		assert.Len(t, handler.proxyRequests, 1)
		expectedProxyRequest := proxyRequest{
			Path:   "/v1/version",
			Method: "GET",
		}
		assert.Equal(t, expectedProxyRequest, handler.proxyRequests[0])
	})
}

func TestHandler_ProxyRoute(t *testing.T) {
	t.Run("should proxy request to destination path if it is configured", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockContext := mock.NewMockContext(ctrl)
		client := mock.NewMockhttpClient(ctrl)
		expectedURL := "/version"
		expectedResponse := &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(""))}
		handler := New(true, true, client)
		handler.proxyRequests = proxyRequests{{
			Path:   "/v1/version",
			Method: "GET",
			Proxy: proxy{
				URL:    expectedURL,
				Method: "GET",
			},
		}}
		mockContext.EXPECT().GetURI().Return(&url.URL{Path: "/v1/version"})
		mockContext.EXPECT().GetMethod().Return("GET")
		client.EXPECT().Do(gomock.Any()).Return(expectedResponse, nil)
		mockContext.EXPECT().SendResponse(expectedResponse, expectedURL)

		handler.ProxyRoute(mockContext)
	})

	t.Run("should return 404 if proxy is not configured", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockContext := mock.NewMockContext(ctrl)
		handler := New(true, true, nil)
		mockContext.EXPECT().GetURI().Return(&url.URL{Path: "/v1/version"})
		mockContext.EXPECT().GetMethod().Return("GET")
		mockContext.EXPECT().Status(404)

		handler.ProxyRoute(mockContext)
	})

	t.Run("should return error if http request creation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockContext := mock.NewMockContext(ctrl)
		client := mock.NewMockhttpClient(ctrl)

		handler := New(true, true, client)
		handler.proxyRequests = proxyRequests{{
			Path:   "/v1/version",
			Method: "GET",
			Proxy: proxy{
				URL:    "foo://bar",
				Method: "😁",
			},
		}}
		mockContext.EXPECT().GetURI().Return(&url.URL{Path: "/v1/version"})
		mockContext.EXPECT().GetMethod().Return("GET")
		mockContext.EXPECT().JSON(400, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.EqualValues(t, "error when creating request for foo://bar: net/http: invalid method \"😁\"", data.(gin.H)["error"])
		})

		handler.ProxyRoute(mockContext)
	})

	t.Run("should return error when proxy fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockContext := mock.NewMockContext(ctrl)
		client := mock.NewMockhttpClient(ctrl)

		handler := New(true, true, client)
		handler.proxyRequests = proxyRequests{{
			Path:   "/v1/version",
			Method: "GET",
			Proxy: proxy{
				URL:    "/version",
				Method: "GET",
			},
		}}
		mockContext.EXPECT().GetURI().Return(&url.URL{Path: "/v1/version"})
		mockContext.EXPECT().GetMethod().Return("GET")
		client.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("error making request"))
		mockContext.EXPECT().JSON(400, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.EqualValues(t, "error when making request for /version: GET", data.(gin.H)["error"])
		})

		handler.ProxyRoute(mockContext)
	})
}

func TestHandler_DeleteProxy(t *testing.T) {
	t.Run("should delete the proxy request if present", func(t *testing.T) {
		handler := New(true, true, nil)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockContext := mock.NewMockContext(ctrl)
		handler.proxyRequests = proxyRequests{proxyRequest{
			Path:   "/v1/version",
			Method: "GET",
		}}
		stringReader := strings.NewReader(`
{
 "path": "/v1/version",
 "method": "GET"
}`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		mockContext.EXPECT().JSON(200, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.EqualValues(t, "deleted the proxy config successfully", data.(gin.H)["result"])
		})

		handler.DeleteProxy(mockContext)

		assert.Len(t, handler.proxyRequests, 0)
	})

	t.Run("should not delete the proxy request if not present", func(t *testing.T) {
		handler := New(true, true, nil)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockContext := mock.NewMockContext(ctrl)
		handler.proxyRequests = proxyRequests{proxyRequest{
			Path:   "/v1/version",
			Method: "GET",
		}}
		stringReader := strings.NewReader(`
{
 "path": "/v2/version",
 "method": "GET"
}`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		mockContext.EXPECT().JSON(404, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.EqualValues(t, "proxy config with url /v2/version and GET method is not found", data.(gin.H)["error"])
		})

		handler.DeleteProxy(mockContext)

		assert.Len(t, handler.proxyRequests, 1)
	})
}
