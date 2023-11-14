package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var telemetryDB = make(map[string]TelemetryData)
var listIPs = make(map[int]string)

// Request is either:
//   - Payload Ops (Ground)
//   - Uplink/Downlink (Ground)
func handleScenarioOne(c *gin.Context, r RedirectRequest) {
	// Define the routes for the two modules
	UplinkRoute := "http://uplink-downlink-module/send/"
	PayloadRoute := "http://payload-ops-module/send/"
}

// Request is either:
//   - Payload Ops (Space)
//   - CNDH (Space)
//   - Uplink/Downlink (Space)
//   - Payload Ops (Center)
func handleScenarioTwo(c *gin.Context, r RedirectRequest) {

	// Create a new HTTP request with the verb, URI and
	// data as specified by the RedirectRequest
	req, err := http.NewRequest(r.Verb, r.URI, bytes.NewReader([]byte(r.Data)))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// Send this request using the default HTTP Client
	// and receive a response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Send the response back to the user
	io.Copy(c.Writer, resp.Body)

}

func receive(c *gin.Context) {
	// Attempt to parse the incoming request's JSON into the "data" struct
	var req RedirectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Abort internally stops Gin from contiuing to handle the request
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Parse the IP from the RedirectRequest
	//
	// We do this by:
	//   1. Calling url.ParseRequestURI which is a Golang
	//      library function used to parse URLs
	//   2. Call .Hostname() which returns the hostname of
	//      the URL, this is where the IP is located
	//
	// E.g. "http://10.1.1.1:8080/telemetry/?id=5"
	//   => "10.1.1.1"

	// Attempt to parse the IP from the RedirectRequest
	parsedURL, err := url.ParseRequestURI(req.URI)
	if err != nil {
		// Handle the error, for example, abort the request
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Get the hostname from the parsed URL
	ip := parsedURL.Hostname()

	// Client is what makes the HTTP request
	resp, err := http.DefaultClient.Do(req)

	// Redirect to specific IPS:
	//   1 - Payload Ops (Space)
	//   2 - CNDH (Space)
	//   3 - Uplink/Downlink (Space)
	//   4 - Uplink/Downlink (Ground)
	//   5 - CNDH (Ground) [US]
	//   6 - Payload Ops (Ground)
	//   7 - Payload Ops (Center)
	switch ip {
	case listIPs[4], listIPs[6]:
		handleScenarioOne(c, req)
		return
	case listIPs[1], listIPs[2], listIPs[3], listIPs[7]:
		handleScenarioTwo(c, req)
		return
	}

	// If we're here, we don't have a valid IP address
	//
	// Note: This could mean we were the IP address and it's
	//       not allowed to send a request back to ourselves
	//       so it's fine if we abort
	c.AbortWithStatus(http.StatusBadRequest)
}

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
	uri := fmt.Sprintf("http://%s:8080/status", listIPs[4])

	//Error handling not implemented on purpose because
	// "An error is returned if there were too many redirects or if there was an HTTP protocol error.
	// A non-2xx response doesn't cause an error."
	res, _ := http.Get(uri)
	defer res.Body.Close()

	// make a 2D array of an interface and string
	var body map[string]interface{}
	//json decoder writes to body by address
	json.NewDecoder(res.Body).Decode(&body)

	//returns the status code and the body in a raw format
	c.JSON(res.StatusCode, body)

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
	server.GET("/receive", receive)
	server.Run() // By default, it will start the server on http://localhost:8080

}
