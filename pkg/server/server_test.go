package server_test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/thecasualcoder/dobby/pkg/server"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestHealth(t *testing.T) {
	t.Run("should return 200", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv, true, true)

		response := performRequest(router, "GET", "/health")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"healthy":true}`, response.Body.String())
	})
}

func TestReadiness(t *testing.T) {
	t.Run("should return 200", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv, true, true)

		response := performRequest(router, "GET", "/readiness")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"ready":true}`, response.Body.String())
	})
}

func TestHealthToggles(t *testing.T) {
	t.Run("should return 500 when sick and 200 when perfect", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv, true, true)

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

		server.Bind(router, srv, true, true)

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

		server.Bind(router, srv, true, true)

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"version":"1.0.0-dev"}`, response.Body.String())
	})

	t.Run("should return 200 with given version", func(t *testing.T) {
		router := gin.Default()
		existingVersion := os.Getenv("VERSION")
		_ = os.Setenv("VERSION", "v1")
		defer func() {
			_ = os.Setenv("VERSION", existingVersion)
		}()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv, true, true)

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"version":"v1"}`, response.Body.String())
	})

	t.Run("should return 500 if service is unhealthy", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv, true, true)

		// make health sick
		performRequest(router, "PUT", "/control/health/sick")

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, `{"error":"application is not healthy"}`, response.Body.String())
	})

	t.Run("should mark service as not healthy till n seconds", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv, true, true)

		// make service not ready
		resetInSeconds := 1
		performRequest(router, "PUT", "/control/health/sick?resetInSeconds="+strconv.Itoa(resetInSeconds))

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, `{"error":"application is not healthy"}`, response.Body.String())

		time.Sleep(time.Duration(resetInSeconds)*time.Second + time.Millisecond)

		response = performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"version":"1.0.0-dev"}`, response.Body.String())
	})

	t.Run("should return 503 if service is not ready", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv, true, true)

		// make service not ready
		performRequest(router, "PUT", "/control/ready/sick")

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusServiceUnavailable, response.Code)
		assert.Equal(t, `{"error":"application is not ready"}`, response.Body.String())
	})

	t.Run("should mark service as not ready till n seconds", func(t *testing.T) {
		router := gin.Default()
		srv := httptest.NewServer(router).Config

		server.Bind(router, srv, true, true)

		// make service not ready
		resetInSeconds := 1
		performRequest(router, "PUT", "/control/ready/sick?resetInSeconds="+strconv.Itoa(resetInSeconds))

		response := performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusServiceUnavailable, response.Code)
		assert.Equal(t, `{"error":"application is not ready"}`, response.Body.String())

		time.Sleep(time.Duration(resetInSeconds)*time.Second + time.Millisecond)

		response = performRequest(router, "GET", "/version")
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, `{"version":"1.0.0-dev"}`, response.Body.String())
	})
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	responseWriter := httptest.NewRecorder()
	r.ServeHTTP(responseWriter, req)
	return responseWriter
}
