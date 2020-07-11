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
}
