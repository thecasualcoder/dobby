package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock "github.com/thecasualcoder/dobby/internal/mock/handler"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestHandler_Call(t *testing.T) {
	t.Run("should return status code alone when there is no body from upstream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mock.NewMockhttpClient(ctrl)
		mockContext := mock.NewMockContext(ctrl)
		handler := New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
 "url": "http://localhost:4444/version",
 "method": "GET"
}`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(""))}, nil)
		mockContext.EXPECT().Status(200)

		handler.Call(mockContext)
	})

	t.Run("should parse the body if the response from upstream contains body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mock.NewMockhttpClient(ctrl)
		mockContext := mock.NewMockContext(ctrl)
		handler := New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
 "url": "http://localhost:4444/version",
 "method": "GET"
}`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(`{"version": 1}`))}, nil)
		mockContext.EXPECT().JSON(200, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.EqualValues(t, 1, data.(map[string]interface{})["version"])
		})

		handler.Call(mockContext)
	})

	t.Run("should return parse error if input is not valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mock.NewMockhttpClient(ctrl)
		mockContext := mock.NewMockContext(ctrl)
		handler := New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
 "url": "http://localhost:4444/notvalid"`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		mockContext.EXPECT().JSON(400, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.Equal(t, "error when decoding request: unexpected EOF", data.(gin.H)["error"])
		})

		handler.Call(mockContext)
	})

	t.Run("should return error if there is error in making request to upstream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mock.NewMockhttpClient(ctrl)
		mockContext := mock.NewMockContext(ctrl)
		handler := New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
 "url": "http://localhost:4444/version"
}`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		httpClient.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("error making request"))
		mockContext.EXPECT().JSON(400, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.Equal(t, "error when making request to http://localhost:4444/version: error making request", data.(gin.H)["error"])
		})

		handler.Call(mockContext)
	})

	t.Run("should parse return special characters in string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mock.NewMockhttpClient(ctrl)
		mockContext := mock.NewMockContext(ctrl)
		handler := New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
 "url": "http://localhost:4444/version"
}`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader("⛅️  +33°C"))}, nil)
		mockContext.EXPECT().JSON(200, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.Equal(t, "⛅️  +33°C", data)
		})

		handler.Call(mockContext)
	})
}

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
