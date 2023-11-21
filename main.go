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
	listIPs = make(map[int]string)

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
	//   1 - Payload Ops (Space)
	//   2 - CNDH (Space)
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
			res, _ := http.DefaultClient.Do(req)
			if res != nil {
				body, _ := io.ReadAll(res.Body)
				fmt.Fprintln(writeLog, "Status ", res.StatusCode, ": ", res.Status, "\nMessage: ", string(body))
			} else {
				fmt.Fprintln(writeLog, "Got a 500")
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
	// Attempt to parse the incoming request's JSON into the "data" struct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
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

func status(c *gin.Context) {
	uri := fmt.Sprintf("http://%s:8080/status", listIPs[4])

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
	server.GET("/status", status)
	server.PUT("/receive", receive)
	server.GET("/execute/:script", executeScript)
	server.GET("/update", updateClient)
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
