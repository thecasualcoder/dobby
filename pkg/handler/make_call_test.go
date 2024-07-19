package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock "github.com/thecasualcoder/dobby/internal/mock/handler"
	h "github.com/thecasualcoder/dobby/pkg/handler"
)

func TestHandler_Call(t *testing.T) {
	t.Run("should make call to upstream on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mock.NewMockhttpClient(ctrl)
		mockContext := mock.NewMockContext(ctrl)
		handler := h.New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
  "url": "http://localhost:4444/version",
  "method": "GET"
}`)
		expectedResponse := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"version": 1}`))}
		mockContext.EXPECT().GetRequestBody().Return(io.NopCloser(stringReader))
		httpClient.EXPECT().Do(gomock.Any()).Return(expectedResponse, nil)
		mockContext.EXPECT().SendResponse(expectedResponse, "http://localhost:4444/version")

		handler.Call(mockContext)
	})

	t.Run("should return parse error if input is not valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mock.NewMockhttpClient(ctrl)
		mockContext := mock.NewMockContext(ctrl)
		handler := h.New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
  "url": "http://localhost:4444/notvalid"`)
		mockContext.EXPECT().GetRequestBody().Return(io.NopCloser(stringReader))
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
		handler := h.New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
  "url": "http://localhost:4444/version"
}`)
		mockContext.EXPECT().GetRequestBody().Return(io.NopCloser(stringReader))
		httpClient.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("error making request"))
		mockContext.EXPECT().JSON(400, gomock.Any()).Do(func(_ int, data interface{}) {
			assert.Equal(t, "error when making request to http://localhost:4444/version: error making request", data.(gin.H)["error"])
		})

		handler.Call(mockContext)
	})
}
