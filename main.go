package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var telemetryDB = make(map[string]TelemetryData)
var listIPs = make(map[int]string)

func putTelemetry(c *gin.Context) {
	id := c.Query("id")    // Extract the ID from the URL path
	var data TelemetryData // Create an empty TelemetryData struct

	// Attempt to parse the incoming request's JSON into the "data" struct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	telemetryDB[id] = data // Store the parsed data in our mock database
	c.JSON(200, data)      // Respond with a 200 status and the stored data
	c.JSON(200, gin.H{"message": "Data saved successfully!"})
}
func getTelemetry(c *gin.Context) {
	id := c.Query("id") // Extract Id from URL path

	if data, ok := telemetryDB[id]; ok {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"coordsX":      data.Coordinates.X,
			"coordsY":      data.Coordinates.Y,
			"coordsZ":      data.Coordinates.Z,
			"temp":         data.Temp,
			"pitch":        data.Rotations.P,
			"yaw":          data.Rotations.Y,
			"roll":         data.Rotations.R,
			"PayloadPower": data.Status.PayloadPower,
			"dataWaiting":  data.Status.DataWaiting,
			"chargeStatus": data.Status.ChargeStatus,
			"voltage":      data.Status.Voltage,
		})
	} else {
		c.JSON(404, gin.H{"error": "data not found!"}) //return 404 if no data
	}
}

func serveFiles(c *gin.Context, contenttype string, path string) {
	filename := c.Param("name")
	filename = path + filename
	_, err := os.Open(filename)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	} else {
		c.Header("Content-Type", contenttype)
		c.File(filename)
	}
}

func serveScripts(c *gin.Context) {
	serveFiles(c, "text/javascript", "./UI/scripts/")
}

func serveCSS(c *gin.Context) {
	serveFiles(c, "text/css", "./UI/styles/")
}

func status(c *gin.Context) {
	uri := fmt.Sprintf("http://%s:8080/test", listIPs[5])
	//I'll have to implement all the error handling
	res, _ := http.Get(uri)
	defer res.Body.Close()
	var body map[string]interface{}
	json.NewDecoder(res.Body).Decode(&body)
	c.JSON(res.StatusCode, body)

}

func test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "hey!"})
}

func readIPCFG() {
	f, err := os.Open("ip.cfg")
	if err != nil {
		fmt.Println("Cannot read ip.cfg.")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		inData := strings.Split(scanner.Text(), ",")
		id, _ := strconv.Atoi(inData[1])
		listIPs[id] = inData[0]
	}
}

func main() {
	readIPCFG()
	server := gin.Default()
	server.LoadHTMLFiles("UI/index.tmpl")
	server.GET("/scripts/:name", serveScripts)
	server.GET("/styles/:name", serveCSS)
	server.PUT("/telemetry/", putTelemetry)
	server.GET("/telemetry/", getTelemetry)
	server.GET("/status", status)
	server.GET("/test", test)
	server.Run() // By default, it will start the server on http://localhost:8080

}
