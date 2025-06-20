package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jmpsec/stanza-c2/pkg/agents"
	"github.com/jmpsec/stanza-c2/pkg/types"
)

// ContentType for header key
const ContentType string = "Content-Type"

// JSONApplicationUTF8 for Content-Type headers
const JSONApplicationUTF8 string = "application/json; charset=UTF-8"

// Standard function to send HTTP requests
func sendHTTPRequest(reqType, url string, params io.Reader, headers map[string]string) (int, []byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(reqType, url, params)
	if err != nil {
		return 0, []byte("Cound not prepare request"), err
	}
	// Prepare headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return 0, []byte("Error sending request"), err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, []byte("Can not read response"), err
	}

	return resp.StatusCode, bodyBytes, nil
}

// Helper to convert a StzRegistrationRequest to an Agent
func agentToRegister(req types.StzRegistrationRequest, ipsrc string) agents.Agent {
	return agents.Agent{
		UUID:     req.UUID,
		Username: req.Username,
		Hostname: req.Hostname,
		GOOS:     req.GOOS,
		GOARCH:   req.GOARCH,
		Uname:    req.Uname,
		IPs:      strings.Join(req.IPs, ", "),
		IPsrc:    ipsrc,
		CycleMin: req.CycleMin,
		CycleMax: req.CycleMax,
		CycleNow: req.CycleMax,
		Active:   true,
	}
}

// Helper to send a serialized JSON response
func httpResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set(ContentType, JSONApplicationUTF8)
	content, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error serializing response: %v", err)
		content = []byte("error serializing response")
	}
	w.WriteHeader(code)
	_, _ = w.Write(content)
}

// Helper to return an error
func httpRegisterResponse(w http.ResponseWriter, resp string) {
	res := types.StzRegistrationResponse{
		Response: resp,
	}
	httpResponse(w, http.StatusOK, res)
}
