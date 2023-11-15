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

type RedirectRequest struct {
	Verb string `json:"verb"`
	URI  string `json:"uri"`
	Data string `json:"data"`
}

var listIPs = make(map[int]string)

func parseScript(scriptName string) (map[int]RedirectRequest, error) {
	f, err := os.Open("scripts/" + scriptName + ".txt")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	allRequests := make(map[int]RedirectRequest)
	count := 0
	i := 0
	scanner := bufio.NewScanner(f)
	var temp RedirectRequest
	for scanner.Scan() {
		input := scanner.Text()
		if input == "STOP" {
			if i == 2 {
				temp.Data = ""
			}
			allRequests[count] = temp
			count++
			i = 0
		} else {
			switch i {
			case 0:
				temp.Verb = input
			case 1:
				indexStr := input[strings.Index(input, "[")+1 : strings.Index(input, "]")]
				index, _ := strconv.Atoi(indexStr)
				temp.URI = strings.Replace(input, "["+indexStr+"]", listIPs[index], 1)
			case 2:
				temp.Data = input
			}
			i++
		}
	}
	return allRequests, nil
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
		req, err := http.NewRequest(temp.Verb, temp.URI, strings.NewReader(temp.Data))
		if err != nil {
			fmt.Fprintln(writeLog, "Failed to make request ", temp.URI, " ", temp.URI)
		} else {
			if temp.Data != "" {
				req.Header.Set("content-type", "application/json")
			}
			res, _ := http.DefaultClient.Do(req)
			if res != nil {
				fmt.Fprintln(writeLog, "Status ", res.StatusCode, ": ", res.Status, "\nMessage: ", res.Body)
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

func readJSONFromFile() (TelemetryData, error) {
	var data TelemetryData
	filename := "telemetry.json"

	file, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()

	// Decode JSON data from the file into the telemetryData variable
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func writeJSONToFile(data TelemetryData) error {
	filename := "telemetry.json"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert TelemetryData to JSON
	dataJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	// Write JSON data to the file
	_, err = file.Write(dataJSON)
	if err != nil {
		return err
	}

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
	uri := fmt.Sprintf("http://%s:8080/status/", listIPs[4])
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
	server.PUT("/telemetry/", putTelemetry)
	server.GET("/telemetry/", getTelemetry)
	server.GET("/status", status)
	server.GET("/execute/:script", executeScript)
	return server
}

func main() {
	temp, err := readIPCFG("ip.cfg")
	if err == nil {
		listIPs = temp
		server := setupServer()
		server.Run() // By default, it will start the server on http://localhost:8080
	} else {
		fmt.Println("Cannot read ip.cfg.")
	}

}
