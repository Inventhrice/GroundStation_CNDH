package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	listIPs      = make(map[int]string)
	clientList   = make(map[chan string]bool)
	newClient    = make(chan chan string)
	closedClient = make(chan chan string)
)

/*
Example request.

curl --header "Content-Type: application/json" \
     --request POST \
     --data '{"verb":"GET","uri":"10.1.1.1/telemetry","data":"example"}' \
     "http://localhost:8080/receive?ID=1"
*/

func receive(c *gin.Context) {
	// Attempt to parse the incoming request's JSON into the "data" struct
	var req RedirectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Abort internally stops Gin from contiuing to handle the request
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Parse the IP from the RedirectRequest
	req.URI = strings.Trim(req.URI, "http://")
	parts := strings.Split(req.URI, "/")
	if len(parts) != 2 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ip := parts[0]

	// Parse the ID query parameter
	stringID := c.Query("ID")
	//if stringID == "" {
	//	c.AbortWithStatus(http.StatusBadRequest)
	//	return
	//}
	sourceID, err := strconv.Atoi(stringID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if sourceID < 1 || sourceID > 7 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Redirect to specific IPS:
	//   1 - CNDH (Space)
	//   2 - Payload Ops (Space)
	//   3 - Uplink/Downlink (Space)
	//   4 - Uplink/Downlink (Ground)
	//   5 - CNDH (Ground) [US]
	//   6 - Payload Ops (Ground)
	//   7 - Payload Ops (Center)
	switch ip {
	case listIPs[4], listIPs[5], listIPs[6]:
		sendRedirectRequest(c, req.Verb, req.URI, []byte(req.Data))
		return
	case listIPs[1], listIPs[2], listIPs[3]:
		// TODO add source ID?
		// uri := listIPs[4] + "/send?ID=" + stringID

		// Create the URI
		uri := listIPs[4] + "/send"

		// Recreate the request data
		data, err := json.Marshal(req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Send the request
		sendRedirectRequest(c, "POST", uri, data)
		return
	case listIPs[7]:
		// Create the URI
		uri := listIPs[6] + "/images"

		// Recreate the request data
		data, err := json.Marshal(req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		sendRedirectRequest(c, "POST", uri, data)
		return
	}
	// If we're here, we don't have a valid IP address
	//
	// Note: This could mean we were the IP address and it's
	//       not allowed to send a request back to ourselves
	//       so it's fine if we abort
	c.AbortWithStatus(http.StatusBadRequest)
}

func executeScript(c *gin.Context) {

	scriptName := c.Param("script")

	writeLog, err := os.Create("scriptOutput.log")
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	defer writeLog.Close()

	allRequests, err := parseScript(scriptName)
	if err != nil {
		writeLog.WriteString(err.Error())
		c.AbortWithStatus(400)
		return
	}

	for i := 0; i < len(allRequests); i++ {
		temp := allRequests[i]

		var req *http.Request
		if temp.Data != "" {
			req, err = http.NewRequest(temp.Verb, temp.URI, bytes.NewBufferString(temp.Data))
			req.Header.Set("content-type", "application/json")
		} else {
			req, err = http.NewRequest(temp.Verb, temp.URI, nil)
		}

		if err != nil {
			fmt.Fprintln(writeLog, "Failed to make request ", temp.URI, " ", temp.URI)
		} else {

			res, err := http.DefaultClient.Do(req)
			if res != nil {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					fmt.Fprintln(writeLog, "Got an error: ", err.Error())
				} else {
					fmt.Fprintln(writeLog, "Status ", res.StatusCode, ": ", res.Status, "\nMessage: ", string(body))
				}

			} else {
				fmt.Fprintln(writeLog, "Got an error: ", err.Error())
			}

		}

	}
	c.Status(200)
}

func getRoot(c *gin.Context) { // Root route reads from json file and puts the data into the html (tmpl) file for display
	c.JSON(200, gin.H{"message": "Server is running"})
}

func putTelemetry(c *gin.Context) {
	//	id := c.Query("id")    // Extract the ID from the URL path (Not currently used)
	var data TelemetryData // Create an empty TelemetryData struct
	body, err := io.ReadAll(c.Request.Body)
	// Attempt to parse the incoming request's JSON into the "data" struct
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&data)
	if err != nil {
		body = body[1 : len(body)-1]
		fmt.Println(string(body))
		err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&data)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	writeErr := writeJSONToFile(data) // Write new data to JSON
	if writeErr != nil {
		c.JSON(400, gin.H{"error": writeErr.Error()})
	} else {
		c.JSON(200, gin.H{"message": "Data saved successfully!"})
	}

	dataJSON, err := json.Marshal(data)
	if err == nil {
		for client := range clientList {
			client <- string(dataJSON)
		}
	}
}

func setTelemetry(c *gin.Context) {
	id := c.Query("id")    // Extract the ID from the URL path
	var data TelemetryData // Create an empty TelemetryData struct

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Read existing data from the file
	existingData, err := readJSONFromFile()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newData := existingData

	if id == "1" { // Only change the data that has been changed
		newData.Coordinates = data.Coordinates
		newData.Rotations = data.Rotations
	} else if id == "2" {
		newData.Coordinates = data.Coordinates
	} else if id == "3" {
		newData.Rotations = data.Rotations
	}

	destAddress := listIPs[1] // Ip for Space CNDH

	moreData := map[string]interface{}{ // Add on the necessary information (Its possible verb isn't needed here?)
		"verb": "PUT",
		"uri":  "http://" + destAddress + ":8080/telemetry/?id=5",
	}

	combinedData := map[string]interface{}{ // Combine the new json data with the telemetry data
		"verb": moreData["verb"],
		"uri":  moreData["uri"],
		"data": newData,
	}

	sendAddress := listIPs[4] // Ip for Ground Uplink/Downlink

	respCode, sendErr := sendTelemetry(c, combinedData, sendAddress) // Function to send the data away
	if sendErr != nil {
		c.JSON(400, gin.H{"error": sendErr.Error()})

	} else {
		writeErr := writeJSONToFile(newData) // Write new json data to file if command went through
		if writeErr != nil {
			c.JSON(400, gin.H{"error": "Data was sent successfully but not saved locally"})
			return
		} else {
			c.JSON(respCode, gin.H{"message": "Successfully saved data and sent command"}) // Should be 200 if everything went properly within the sendTelemetry function (400 if timeout)

			dataJSON, err := json.Marshal(newData)
			if err == nil {
				for client := range clientList {
					client <- string(dataJSON)
				}
			}
		}
	}
}

func getTelemetry(c *gin.Context) {
	//	id := c.Query("id")

	data, err := readJSONFromFile() // Load json data into data variable
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{ // Write json data to html page
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

func updateClient(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Connection", "keep-alive")
	c.Header("Cache-Control", "no-cache")

	client := make(chan string, 1)
	newClient <- client
	defer func() {
		closedClient <- client
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case data, ok := <-client:
			if !ok {
				return false
			}
			c.SSEvent("message", data)
			return true

		case <-c.Request.Context().Done():
			return false
		}
	})
}

func serveScripts(c *gin.Context) {
	serveFiles(c, "text/javascript", "./UI/scripts/")
}

func serveCSS(c *gin.Context) {
	serveFiles(c, "text/css", "./UI/styles/")
}

func requestTelemetry(c *gin.Context) {
	uri := fmt.Sprintf("http://%s:8080/send/", listIPs[4])

	// Create JSON
	json := "{\"verb\":\"GET\",\"uri\":\"http://" + listIPs[1] + ":8080/send/\"}"
	body := strings.NewReader(json)

	res, err := http.NewRequest("POST", uri, body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer res.Body.Close()
	c.JSON(http.StatusOK, gin.H{"message": "Telemetry requested"})
}

func status(c *gin.Context) {
	uri := fmt.Sprintf("http://%s:8080/status/", listIPs[4])

	// Error handling not implemented on purpose because
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

func readIPCFG(path string) (map[int]string, error) {
	ips := make(map[int]string)
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			inData := strings.Split(scanner.Text(), ",")
			id, _ := strconv.Atoi(inData[1])
			ips[id] = inData[0]
		}
	}
	return ips, err
}

func setupServer() *gin.Engine {
	server := gin.Default()
	server.LoadHTMLFiles("UI/index.tmpl")
	server.GET("/", getRoot)
	server.GET("/scripts/:name", serveScripts)
	server.GET("/styles/:name", serveCSS)
	server.PUT("/telemetry", putTelemetry)
	server.GET("/telemetry", getTelemetry)
	server.PUT("/settelemetry", setTelemetry)
	server.GET("/status", status)
	server.PUT("/receive", receive)
	server.GET("/execute/:script", executeScript)
	server.GET("/update", updateClient)
	server.GET("/requestTelemetry", requestTelemetry)
	return server
}

func main() {
	go manageClientList()
	temp, err := readIPCFG("ip.cfg")
	if err == nil {
		listIPs = temp
		server := setupServer()
		server.Run() // By default, it will start the server on http://localhost:8080
	} else {
		fmt.Println("Cannot read ip.cfg.")
	}
}
