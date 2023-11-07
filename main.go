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

var telemetryData TelemetryData
var listIPs = make(map[int]string)
var gotData = false // Updates upon succesful loading of JSON data from putTelemetry route for check in getTelemetry route

func getRoot(c *gin.Context) { // Root route reads from json file and puts the data into the html (tmpl) file for display
	err := readJSONFromFile()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	data := telemetryData
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
}

func putTelemetry(c *gin.Context) {
	//	id := c.Query("id")    // Extract the ID from the URL path (Not currently used)
	var data TelemetryData // Create an empty TelemetryData struct

	// Attempt to parse the incoming request's JSON into the "data" struct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	telemetryData = data // Store the parsed data in our mock database
	gotData = true
	c.JSON(200, data) // Respond with a 200 status and the stored data

	writeErr := writeJSONToFile() // Write new data to JSON
	if writeErr != nil {
		fmt.Print("Error", writeErr)
		c.JSON(400, gin.H{"error": writeErr.Error()})
	} else {
		c.JSON(200, gin.H{"message": "Data saved successfully!"})
	}
}

func getTelemetry(c *gin.Context) {
	//	id := c.Query("id")
	if data := telemetryData; gotData {
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

		err := writeJSONToFile()
		if err != nil {
			fmt.Println("Error", err)
		}

	} else {
		c.JSON(404, gin.H{"error": "data not found!"}) //return 404 if no data
	}
	gotData = false // Reset gotData so we can check again on next request
}

func readJSONFromFile() error {
	filename := "telemetry.json"

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode JSON data from the file into the telemetryData variable
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&telemetryData)
	if err != nil {
		return err
	}

	return nil
}

func writeJSONToFile() error {
	filename := "telemetry.json"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert TelemetryData to JSON
	dataJSON, err := json.MarshalIndent(telemetryData, "", "    ")
	if err != nil {
		return err
	}

	// Write JSON data to the file
	_, err = file.Write(dataJSON)
	if err != nil {
		return err
	}

	fmt.Println("JSON data written to", filename)
	return nil
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
	server.GET("/", getRoot)
	server.GET("/scripts/:name", serveScripts)
	server.GET("/styles/:name", serveCSS)
	server.PUT("/telemetry/", putTelemetry)
	server.GET("/telemetry/", getTelemetry)
	server.GET("/status", status)
	server.Run() // By default, it will start the server on http://localhost:8080

}
