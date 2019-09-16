package server_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/dobby/pkg/server"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealth(t *testing.T) {
	t.Run("should return 200", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv)

		response := performRequest(router, "GET", "/health")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"healthy":true}`, response.Body.String())
	})
}

func TestReadiness(t *testing.T) {
	t.Run("should return 200", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv)

		response := performRequest(router, "GET", "/readiness")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"ready":true}`, response.Body.String())
	})
}

func TestHealthToggles(t *testing.T) {
	t.Run("should return 500 when sick and 200 when perfect", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv)

		response := performRequest(router, "PUT", "/control/health/sick")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"status":"success"}`, response.Body.String())

		response = performRequest(router, "GET", "/health")
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, `{"healthy":false}`, response.Body.String())

		response = performRequest(router, "PUT", "/control/health/perfect")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"status":"success"}`, response.Body.String())

		response = performRequest(router, "GET", "/health")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"healthy":true}`, response.Body.String())
	})
}

func TestReadinessToggles(t *testing.T) {
	t.Run("should return 500 when sick and 200 when perfect", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv)

		response := performRequest(router, "PUT", "/control/ready/sick")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"status":"success"}`, response.Body.String())

		response = performRequest(router, "GET", "/readiness")
		assert.Equal(t, http.StatusServiceUnavailable, response.Code)
		assert.Equal(t, `{"ready":false}`, response.Body.String())

		response = performRequest(router, "PUT", "/control/ready/perfect")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"status":"success"}`, response.Body.String())

		response = performRequest(router, "GET", "/readiness")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"ready":true}`, response.Body.String())
	})
}

func TestVersion(t *testing.T) {
	t.Run("should return 200 with version", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv)

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"version":"1.0.0-dev"}`, response.Body.String())
	})

	t.Run("should return 200 with given version", func(t *testing.T) {
		router := gin.Default()
		_ = os.Setenv("VERSION", "v1")
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv)

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"version":"v1"}`, response.Body.String())
	})

	t.Run("should return 500 if service is unhealthy", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv)

		// make health sick
		performRequest(router, "PUT", "/control/health/sick")

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, `{"error":"application is not healthy"}`, response.Body.String())

		// make health perfect again
		performRequest(router, "PUT", "/control/health/perfect")
	})

	t.Run("should return 503 if service is not ready", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv)

		// make service not ready
		performRequest(router, "PUT", "/control/ready/sick")

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusServiceUnavailable, response.Code)
		assert.Equal(t, `{"error":"application is not ready"}`, response.Body.String())

		// make readiness perfect again
		performRequest(router, "PUT", "/control/ready/perfect")
	})
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	responseWriter := httptest.NewRecorder()
	r.ServeHTTP(responseWriter, req)
	return responseWriter
}
