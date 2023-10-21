package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Coordinates struct {
	X string `json:"x"`
	Y string `json:"y"`
	Z string `json:"z"`
}

type Rotations struct {
	P string `json:"p"`
	Y string `json:"y"`
	R string `json:"r"`
}

type Status struct {
	PayloadPower string `json:"payloadPower"`
	DataWaiting  string `json:"dataWaiting"`
	ChargeStatus string `json:"chargeStatus"`
	Voltage      string `json:"voltage"`
}

type TelemetryData struct {
	Coordinates Coordinates `json:"coordinates"`
	Rotations   Rotations   `json:"rotations"`
	Fuel        string      `json:"fuel"`
	Temp        string      `json:"temp"`
	Status      Status      `json:"status"`
}

var telemetryDB = make(map[string]TelemetryData)
var listIPs = make(map[int]string)

func main() {
	f, err := os.Open("ip.cfg")
	if err != nil {
		fmt.Println("Cannot read ip.cfg.")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		inData := strings.SplitAfter(scanner.Text(), ",")
		id, _ := strconv.Atoi(inData[1])
		listIPs[id] = inData[0]
	}

	server := gin.Default()
	server.LoadHTMLFiles("UI/index.tmpl")

	// Default route - FINISH THIS
	server.GET("/", func(c *gin.Context) {
		fmt.Println(c.FullPath())
		c.JSON(200, gin.H{"message": "Data saved successfully!"})
	})

	server.PUT("/telemetry/:id", func(c *gin.Context) {
		id := c.Param("id")    // Extract the ID from the URL path
		var data TelemetryData // Create an empty TelemetryData struct

		// Attempt to parse the incoming request's JSON into the "data" struct
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		telemetryDB[id] = data // Store the parsed data in our mock database
		c.JSON(200, data)      // Respond with a 200 status and the stored data
		c.JSON(200, gin.H{"message": "Data saved successfully!"})
	})

	server.GET("/telemetry/:id", func(c *gin.Context) {
		id := c.Param("id") // Extract Id from URL path

		if data, ok := telemetryDB[id]; ok {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"coordsX": data.Coordinates.X,
				"coordsY": data.Coordinates.Y,
				"coordsZ": data.Coordinates.Z,
				"temp":    data.Temp,
				"pitch":   data.Rotations.P,
				"yaw":     data.Rotations.Y,
				"roll":    data.Rotations.R,
			})
			//c.JSON(200, data)
		} else {
			c.JSON(404, gin.H{"error": "data not found!"}) //return 404 if no data
		}
	})

	server.Run() // By default, it will start the server on http://localhost:8080
}
