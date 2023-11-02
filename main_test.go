package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestServer() *gin.Engine {
	server := gin.Default()
	return server
}

func Test_GetTelemetry(t *testing.T) {
	expected := `{"error":"data not found!"}`
	server := gin.Default()
	server.GET("/telemetry/0", getTelemetry)
	req, _ := http.NewRequest("GET", "/telemetry/0", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	actual, _ := io.ReadAll(w.Body)
	assert.Equal(t, expected, string(actual))
	assert.Equal(t, http.StatusNotFound, w.Code)
}
