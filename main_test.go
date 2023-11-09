package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestServer() *gin.Engine {
	server := gin.Default()
	return server
}

func Test01_GetTelemetry_NoData(t *testing.T) {
	expected := `{"error":"data not found!"}`
	server := gin.Default()
	server.GET("/telemetry/", getTelemetry)
	req, _ := http.NewRequest("GET", "/telemetry/?id=0", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	actual, _ := io.ReadAll(w.Body)
	assert.Equal(t, expected, string(actual))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test02_GetTelemetry_WithData(t *testing.T) {
	// Create mock data
	coords := Coordinates{X: "1", Y: "1", Z: "1"}
	rot := Rotations{P: "1", Y: "1", R: "1"}
	stat := Status{
		PayloadPower: "",
		DataWaiting:  "",
		ChargeStatus: "",
		Voltage:      "",
	}
	data := TelemetryData{Coordinates: coords, Rotations: rot, Status: stat}

	// Set the telemetryDB with mock data
	var index = "0"
	telemetryDB[index] = data
	filename := "test/test2.tmpl"
	expectedBody, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening test html file.")
	}
	expectedCode := http.StatusOK
	server := gin.Default()
	server.LoadHTMLFiles("UI/index.tmpl")
	server.GET("/telemetry/", getTelemetry)
	req, _ := http.NewRequest("GET", "/telemetry/?id=0", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	actual, _ := io.ReadAll(w.Body)
    //if err := os.WriteFile(filename, []byte(string(actual)), 0666); err != nil {
    // fmt.Println(err)
    //}
	assert.Equal(t, string(expectedBody), string(actual))
	assert.Equal(t, expectedCode, w.Code)
}

func Test03_PutTelemetry_InvalidInput(t *testing.T) {
    data := "vasdvasdvasdafsadvasdv"
    expectedCode := http.StatusBadRequest
    server := gin.Default()
    server.PUT("/telemetry/", putTelemetry)
    req, _ := http.NewRequest("PUT", "/telemetry/?id=0,data="+data, nil)
    w := httptest.NewRecorder()
    server.ServeHTTP(w, req)
    assert.Equal(t, expectedCode, w.Code)
}

func Test04_Scripts(t *testing.T){
	expected := "testTextfor js file"
	
	server := gin.Default()
	server.GET("/scripts/:name", serveScripts)
	
}

func Test05_Styles(t *testing.T){
	expected := "testTextfor css file"

	server := gin.Default()
	server.GET("/styles/:name", serveCSS)
}

func Test06_HTML(t *testing.T){
	
}

func Test07_Root(t *testing.T){
	expected := http.StatusOK
	expectedMsg := "Server is running"

	server := gin.Default()
	server.GET("/", getRoot)
}

func test08_readIPCFG(t *testing.T){
	
}