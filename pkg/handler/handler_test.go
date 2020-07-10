package handler

import (
	"github.com/golang/mock/gomock"
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

	t.Run("should return status code alone when there is no body from upstream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mock.NewMockhttpClient(ctrl)
		mockContext := mock.NewMockContext(ctrl)
		handler := New(true, true, httpClient)

		stringReader := strings.NewReader(`
{
 "url": "http://localhost:4444/return/500",
 "method": "GET"
}`)
		mockContext.EXPECT().GetRequestBody().Return(ioutil.NopCloser(stringReader))
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{StatusCode: 500, Body: ioutil.NopCloser(strings.NewReader(""))}, nil)
		mockContext.EXPECT().Status(500)

		handler.Call(mockContext)
	})
}
