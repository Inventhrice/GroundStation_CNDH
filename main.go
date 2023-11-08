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

// Request represents the JSON data structure for incoming requests.
type Request struct {
	Destination string `json:"destination"`
	Verb        string `json:"verb"`
	IP          string `json:"ip"`
	Route       string `json:"route"`
}

func receive(c *gin.Context) {
	jsonData := c.Query("data") // data parsed from query parameter
	var request Request         // the struct we made to parse the JSON
	if err := json.Unmarshal([]byte(jsonData), &request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON data"}) //cannot parse
		return
	}
	destination := request.Destination
	verb := request.Verb
	ip := request.IP
	route := request.Route

	switch destination {
	case "GroundPayloadOps":
		if verb == "GET" { // Handle GET request for GroundPayloadOps
			resp, err := http.Get("http://" + ip + route)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to make GET request"})
				return
			}
			defer resp.Body.Close() // cleanup and release HTTP request 

			c.JSON(200, gin.H{"message": "GET request processed successfully"})
		} else if verb == "PUT" { // Handle PUT request for GroundPayloadOps

			c.JSON(200, gin.H{"message": "PUT request processed successfully"})
		} else {
			c.JSON(400, gin.H{"error": "Unsupported HTTP verb"})
		}
	case "UplinkDownlink":
		if verb == "GET" { // Handle GET request for UplinkDownlink
			resp, err := http.Get("http://" + ip + route)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to make GET request"})
				return
			}
			defer resp.Body.Close() // cleanup and release HTTP request 

			c.JSON(200, gin.H{"message": "GET request processed successfully"})
		} else if verb == "PUT" {

			c.JSON(200, gin.H{"message": "PUT request processed successfully"})
		} else {
			c.JSON(400, gin.H{"error": "Unsupported HTTP verb"})
		}

	default://destination not in our scope so we Forward the request
		resp, err := http.NewRequest(verb, "http://"+ip+route, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create the request"})
			return
		}
	
		client := &http.Client{}
		resp, err = client.Do(resp)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to make the request"})
			return
		}
		defer resp.Body.Close() // cleanup and release HTTP request 
	
		c.JSON(200, gin.H{"message": "Request processed successfully"})

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
	server.GET("/recieve")
	server.Run() // By default, it will start the server on http://localhost:8080

}
