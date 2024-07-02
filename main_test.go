package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_PingHandler(t *testing.T) {
	w := httptest.NewRecorder()
	pingHandler(http.ResponseWriter(w), time.Duration(0))
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_PingHandler_WithDuration(t *testing.T) {
	w := httptest.NewRecorder()
	now := time.Now()
	workDuration := time.Duration(10 * time.Millisecond)
	pingHandler(http.ResponseWriter(w), workDuration)
	resp := w.Result()
	elapsed := time.Since(now)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Greater(t, elapsed, workDuration)
}

func Test_VersionHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	versionHandler(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, versionString, w.Body.String())
}

func Test_HealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_HealthHandler_WithHealthyParamFalse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health?healthy=false", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_StatusHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	statusHandler(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_StatusHandler_Status418(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/status?status=418", nil)
	w := httptest.NewRecorder()
	statusHandler(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}

func Test_EchoHandler(t *testing.T) {
	// create an io.Reader with JSON payload { "message": "hello world"}
	jsonPayload := []byte(`{"message": "hello world"}`)
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBuffer(jsonPayload))
	w := httptest.NewRecorder()
	echoHandler(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, string(jsonPayload), w.Body.String())
}
