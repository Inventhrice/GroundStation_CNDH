package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := setupServer()
	return server
}

func Test04_Scripts_Valid(t *testing.T) {
	expected := "testTextfor js file"

	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/scripts/test.js", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Body.String())
}

func Test05_Styles_Valid(t *testing.T) {
	expected := "testTextfor css file"

	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/styles/test.css", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Body.String())
}

func Test04_Scripts_Invalid(t *testing.T) {
	expected := "{\"error\":\"open ./UI/scripts/NOTFOUND.js: The system cannot find the file specified.\"}"

	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/scripts/NOTFOUND.js", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Body.String())
}

func Test05_Styles_Invalid(t *testing.T) {
	expected := "{\"error\":\"open ./UI/styles/NOTFOUND.css: The system cannot find the file specified.\"}"

	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/styles/NOTFOUND.css", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Body.String())
}

// func Test06_HTML(t *testing.T){

// }

func Test07_Root(t *testing.T) {
	expected := http.StatusOK
	expectedMsg := "{\"message\":\"Server is running\"}"

	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Code)
	assert.Equal(t, expectedMsg, w.Body.String())
}

// func test08_readIPCFG(t *testing.T){

// }
