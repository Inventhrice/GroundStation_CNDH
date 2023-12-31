package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func sendRedirectRequest(c *gin.Context, verb string, uri string, data []byte) {
	// Creates the request
	reader := bytes.NewReader(data)

	r, err := http.NewRequest(verb, "http://"+uri, reader)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Send the request
	client := &http.Client{
		// Manually add a timeout of 10 seconds
		// because the default client does not
		// contain one
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(r)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Replies with status code of request
	c.Status(resp.StatusCode)
	return
}

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

func manageClientList() {
	for {
		select {
		case client := <-newClient:
			clientList[client] = true
		case client := <-closedClient:
			delete(clientList, client)
			close(client)
		}
	}
}

func sendTelemetry(c *gin.Context, combinedData map[string]interface{}, ipAddress string) (int, error) {

	// Convert TelemetryData to JSON
	dataJSON, marshErr := json.Marshal(combinedData)
	if marshErr != nil {
		return 0, marshErr
	}

	// Send a PUT request to the specified IP address
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := "http://" + ipAddress + ":8080/send/" // Update the URL endpoint accordingly

	req, dataErr := http.NewRequest("POST", url, bytes.NewBuffer(dataJSON)) // Create the request
	if dataErr != nil {
		return 0, dataErr
	}

	req.Header.Set("Content-Type", "application/json")

	resp, respErr := client.Do(req) // Collect a response
	if respErr != nil {
		return 0, respErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // Should be 200 if everything went well
		return resp.StatusCode, respErr
	}

	return resp.StatusCode, nil
}
