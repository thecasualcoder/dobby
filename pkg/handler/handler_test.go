package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	h "github.com/thecasualcoder/dobby/pkg/handler"
)

func TestDefaultContext_SendResponse(t *testing.T) {
	url := "http://localhost:4444/version"

	t.Run("should return status code alone when there is no body from upstream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		recorder := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(recorder)
		context := h.NewDefaultContext(ginContext)
		httpResponse := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}

		context.SendResponse(httpResponse, url)

		assert.Equal(t, 200, recorder.Code)
		assert.Empty(t, recorder.Body)
	})

	t.Run("should parse the body if the response from upstream contains body", func(t *testing.T) {
		t.Run("parse object response", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			context := h.NewDefaultContext(ginContext)
			httpResponse := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"version": 1}`))}

			context.SendResponse(httpResponse, url)

			assert.Equal(t, 200, recorder.Code)
			assert.NotEmpty(t, recorder.Body)
			actualResponse := make(gin.H)
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &actualResponse))
			assert.EqualValues(t, 1, actualResponse["version"])
		})

		t.Run("parse array response", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			context := h.NewDefaultContext(ginContext)
			httpResponse := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`[{"version": 1}]`))}

			context.SendResponse(httpResponse, url)

			assert.Equal(t, 200, recorder.Code)
			assert.NotEmpty(t, recorder.Body)
			var actualResponse []gin.H
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &actualResponse))
			assert.EqualValues(t, 1, actualResponse[0]["version"])
		})

		for _, expectedResponse := range []string{"someData", "⛅️  +33°C"} {
			t.Run(fmt.Sprintf("parse string response %s", expectedResponse), func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				recorder := httptest.NewRecorder()
				ginContext, _ := gin.CreateTestContext(recorder)
				context := h.NewDefaultContext(ginContext)
				httpResponse := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(expectedResponse))}

				context.SendResponse(httpResponse, url)

				assert.Equal(t, 200, recorder.Code)
				assert.NotEmpty(t, recorder.Body)
				var actualResponse string
				assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &actualResponse))
				assert.EqualValues(t, expectedResponse, actualResponse)
			})
		}
	})

	t.Run("should return parse error if input is not valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		recorder := httptest.NewRecorder()
		ginContext, _ := gin.CreateTestContext(recorder)
		context := h.NewDefaultContext(ginContext)
		httpResponse := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"url": "http://localhost:4444/notvalid"`))}

		context.SendResponse(httpResponse, url)

		assert.Equal(t, 400, recorder.Code)
		assert.NotEmpty(t, recorder.Body)
		actualResponse := make(gin.H)
		assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &actualResponse))
		assert.EqualValues(t, "error when decoding response from http://localhost:4444/version: unexpected end of JSON input", actualResponse["error"])
	})
}
