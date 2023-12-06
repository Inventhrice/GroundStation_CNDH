package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func makeTestListIP() map[int]string {
	temp := make(map[int]string)
	for i := 1; i <= 7; i++ {
		temp[i] = "localhost"
	}
	return temp
}

func SetupTestServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := setupServer()
	go manageClientList()
	return server
}

func Test01_PutTelemetry_ValidInput(t *testing.T) {
	// Create input data
	coords := Coordinates{X: 1, Y: 1, Z: 1}
	rot := Rotations{P: 1, Y: 1, R: 1}
	stat := Status{PayloadPower: "", DataWaiting: true, ChargeStatus: true, Voltage: 12.5}
	data := TelemetryData{Coordinates: coords, Rotations: rot, Status: stat, Fuel: 80, Temp: 80.4}
	// Convert data into JSON encoded byte array
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(jsonData)
	req, _ := http.NewRequest("PUT", "/telemetry", body)
	w := httptest.NewRecorder()
	server := SetupTestServer()
	server.ServeHTTP(w, req)
	actual, _ := io.ReadAll(w.Body)
	expectedBody := "{\"message\":\"Data saved successfully!\"}"
	expectedCode := http.StatusOK

	actualFilename := "telemetry.json"
	actualJson, err := os.ReadFile(actualFilename)
	if err != nil {
		t.Fatal(err)
	}
	expectedFilename := "FilesForTesting/Test1_Expected.json"
	expectedJson, err := os.ReadFile(expectedFilename)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(expectedBody), string(actual))
	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedJson, actualJson)
}

func Test02_PutTelemetry_InvalidInput(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/telemetry", nil)
	w := httptest.NewRecorder()
	server := SetupTestServer()
	server.ServeHTTP(w, req)

	expectedCode := http.StatusInternalServerError
	assert.Equal(t, expectedCode, w.Code)
}

func Test03_GetTelemetry_Valid(t *testing.T) {
	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/telemetry", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	actual, _ := io.ReadAll(w.Body)

	filename := "FilesForTesting/Test3_Expected.tmpl"
	//if err := os.WriteFile(filename, []byte(string(actual)), 0666); err != nil {
	//    t.Fatal(err)
	//}
	expectedBody, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	expectedCode := http.StatusOK
	assert.Equal(t, string(expectedBody), string(actual))
	assert.Equal(t, expectedCode, w.Code)
}

func Test04_Scripts_Valid(t *testing.T) {
	expected := "testTextfor js file"

	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/scripts/test.js", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Body.String())
}

func Test05_Scripts_Invalid(t *testing.T) {
	expected := "{\"error\":\"open ./UI/scripts/NOTFOUND.js: no such file or directory\"}"
	expectedCode := http.StatusNotFound
	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/scripts/NOTFOUND.js", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Body.String())
	assert.Equal(t, expectedCode, w.Code)
}

func Test06_Styles_Valid(t *testing.T) {
	expected := "testTextfor css file"

	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/styles/test.css", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Body.String())
}

func Test07_Styles_Invalid(t *testing.T) {
	expected := "{\"error\":\"open ./UI/styles/NOTFOUND.css: no such file or directory\"}"
	expectedCode := http.StatusNotFound
	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/styles/NOTFOUND.css", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Body.String())
	assert.Equal(t, expectedCode, w.Code)
}

func Test08_Root(t *testing.T) {
	expected := http.StatusOK
	expectedMsg := "{\"message\":\"Server is running\"}"

	server := SetupTestServer()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, expected, w.Code)
	assert.Equal(t, expectedMsg, w.Body.String())
}

/* func Test09_readIPCFG_Valid(t *testing.T) {
	_, err := readIPCFG("ip.cfg")

	assert.Equal(t, nil, err)
}*/

func Test10_readIPCFG_Invalid(t *testing.T) {

	_, err := readIPCFG("nilpath.cfg")
	assert.Equal(t, "open nilpath.cfg: no such file or directory", err.Error())
}

func Test11_executeScript_Valid(t *testing.T) {
	expectedCode := 200

	listIPs = makeTestListIP()
	server := SetupTestServer()
	req, _ := http.NewRequest("GET", "/execute/Script1", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, expectedCode, w.Code)
}
func Test12_executeScript_InvalidScript(t *testing.T) {
	expectedCode := 400

	server := SetupTestServer()
	req, _ := http.NewRequest("GET", "/execute/NOTFOUNDFILE", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, expectedCode, w.Code)

}
