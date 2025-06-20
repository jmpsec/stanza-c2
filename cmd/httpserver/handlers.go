package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/jmpsec/stanza-c2/pkg/agents"
	"github.com/jmpsec/stanza-c2/pkg/types"
)

// Handle requests with empty responses
func emptyHTTPHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))

	// Send response
	httpResponse(w, http.StatusOK, "Soon...")
}

// Handle all HTTP requests and redirect them
func redirectHTTPHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))
	http.Redirect(w, r, redirectURL, 301)
}

// Handle HTTP registration requests
func registerHTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("STZ: Registration request!\n")
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))
	// Read and decode POST body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		httpRegisterResponse(w, types.StzResponseError)
		return
	}
	var reg types.StzRegistrationRequest
	if err := json.Unmarshal(body, &reg); err != nil {
		log.Println(err)
		httpRegisterResponse(w, types.StzResponseError)
		return
	}
	// If agent exists, update otherwise register
	ip := r.Header.Get("X-Real-IP")
	agent := agentToRegister(reg, ip)
	_, alreadyExists := stzAgents.Exist(agent.UUID, agent.Hostname, ip, agent.Username)
	if alreadyExists {
		if err := stzAgents.Update(agent); err != nil {
			log.Println(err)
			httpRegisterResponse(w, types.StzResponseError)
			return
		}
		log.Printf("STZ_UPDATE Agent!\n")
		// Log agent activity for update
		entry := agents.AgentLog{
			UUID:   agent.UUID,
			IPsrc:  ip,
			Action: types.StzActionUpdate,
			Data:   string(body),
		}
		if err := stzAgents.Log(&entry); err != nil {
			log.Println(err)
		}
	} else {
		if err := stzAgents.Register(&agent); err != nil {
			log.Println(err)
			httpRegisterResponse(w, types.StzResponseError)
			return
		}
		log.Printf("STZ_REGISTER Agent!\n")
		// Log agent activity for register
		entry := agents.AgentLog{
			UUID:   agent.UUID,
			IPsrc:  ip,
			Action: types.StzActionRegister,
			Data:   string(body),
		}
		if err := stzAgents.Log(&entry); err != nil {
			log.Println(err)
			httpRegisterResponse(w, types.StzResponseError)
			return
		}
	}
	// Send response
	httpRegisterResponse(w, types.StzResponseOk)
}

// Handle HTTP beacon requests
func beaconHTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("STZ: Beacon received!\n")
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))
	// Read and decode POST body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var req types.StzBeaconStatus
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println(err)
		return
	}
	// Check if agent exists
	ip := r.Header.Get("X-Real-IP")
	_, e := stzAgents.ExistBeacon(req.UUID, ip)
	if !e {
		register := types.StzBeaconResponse{
			ID:      0,
			Action:  types.StzActionRegister,
			Payload: "",
		}
		httpResponse(w, http.StatusOK, []types.StzBeaconResponse{register})
		return
	}
	// Update current cycle value
	if err := stzAgents.UpdateBeaconCycle(req.UUID, req.CycleNow); err != nil {
		log.Println(err)
		return
	}
	// Logging activity
	entry := agents.AgentLog{
		UUID:   req.UUID,
		IPsrc:  ip,
		Action: types.StzStatusBeacon,
		Data:   string(body),
	}
	if err := stzAgents.Log(&entry); err != nil {
		log.Println(err)
	}
	// Retrieve commands for this UUID
	commands, err := stzCommands.BeaconCommands(req.UUID)
	if err != nil {
		log.Println(err)
	}
	// Send response
	httpResponse(w, http.StatusOK, commands)
}

// Handle HTTP execution requests
func executionHTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("STZ: Execution!\n")
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))

	// Read and decode POST body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var req types.StzExecutionStatus
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println(err)
		return
	}
	// Update the status in DB for that command, if it is valid
	validStatus := map[string]bool{
		types.StzStatusReceived: true,
		types.StzStatusDone:     true,
	}
	// Check if status is valid and update
	if validStatus[req.Status] {
		err := stzCommands.Update(req.ID, req.Status, req.Data)
		if err != nil {
			log.Println(err)
		}
	}
	// Logging activity
	// TODO update ipaddress
	ip := r.Header.Get("X-Real-IP")
	if err := stzAgents.LogCheckin(req.UUID, ip, req.Status, string(body)); err != nil {
		log.Println(err)
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte{})
}

// Handle HTTP callbacks requests
func callbacksHTTPHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("STZ: Callbacks!\n")
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))

	// Read and decode POST body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var req types.StzCallbacksRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println(err)
	}
	// Check valid passcode
	if req.HelloThisIsDog != httpConfig.Token {
		w.Header().Set("Content-Type", JSONApplicationUTF8)
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte{})
		return
	}
	// Get callbacks from DB
	callbacks, err := stzCallbacks.StzCallbacks()
	if err != nil {
		log.Println(err)
		return
	}
	// Logging activity
	ip := r.Header.Get("X-Real-IP")
	entry := agents.AgentLog{
		UUID:   req.UUID,
		IPsrc:  ip,
		Action: types.StzActionCallback,
		Data:   string(body),
	}
	if err := stzAgents.Log(&entry); err != nil {
		log.Println(err)
	}
	// Prepare response
	response, err := json.Marshal(
		types.StzCallbacksResponse{
			Callbacks: callbacks,
		})
	if err != nil {
		log.Println(err)
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}
