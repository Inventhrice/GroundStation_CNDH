package main

import (
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

func Test04_Scripts(t *testing.T){
	expected := "testTextfor js file"
	
	server := gin.Default()
	server.GET("/scripts/:name", serveScripts)

	req, _ := http.NewRequest("GET", "/scripts/test.js", nil)
    w := httptest.NewRecorder()

    server.ServeHTTP(w, req)

    assert.Equal(t, expected, w.Body.String())
}

func Test05_Styles(t *testing.T){
	expected := "testTextfor css file"

	server := gin.Default()
	server.GET("/styles/:name", serveCSS)

	req, _ := http.NewRequest("GET", "/styles/test.css", nil)
    w := httptest.NewRecorder()

    server.ServeHTTP(w, req)

    assert.Equal(t, expected, w.Body.String())
}

// func Test06_HTML(t *testing.T){
	
// }

// func Test07_Root(t *testing.T){
// 	expected := http.StatusOK
// 	expectedMsg := "Server is running"

// 	server := gin.Default()
// 	server.GET("/", getRoot)
// }

// func test08_readIPCFG(t *testing.T){

// }