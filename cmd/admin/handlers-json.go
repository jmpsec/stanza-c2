package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/jmpsec/stanza-c2/pkg/commands"
	"github.com/jmpsec/stanza-c2/pkg/types"
)

// Handler to create new actions in the C2
func commandsActionHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	// Make sure the action is valid
	validActions := map[string]bool{
		types.StzActionSet:      true,
		types.StzActionRegister: true,
		types.StzActionExecute:  true,
		types.StzActionGet:      true,
		types.StzActionPut:      true,
		types.StzActionDelete:   true,
		types.StzActionLock:     true,
		types.StzActionSleep:    true,
		types.StzActionExit:     true,
	}
	action := r.PathValue("action")
	if action == "" {
		w.Header().Set("Content-Type", JSONApplicationUTF8)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("NO ACTION"))
		return
	}
	if !validActions[strings.ToUpper(action)] {
		w.Header().Set("Content-Type", JSONApplicationUTF8)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD ACTION"))
		return
	}
	// Read and decode POST body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var newCmd StzNewCommand
	if err := json.Unmarshal(body, &newCmd); err != nil {
		log.Println(err)
	}
	// Action should be the same
	if action == newCmd.Action {
		for _, t := range newCmd.Targets {
			// Create command
			cmd := commands.Command{
				Target:    t,
				Action:    newCmd.Action,
				Payload:   newCmd.Payload,
				Completed: false,
			}
			if err := stzCommands.New(&cmd); err != nil {
				log.Printf("error creating new command %v", err)
				w.Header().Set("Content-Type", JSONApplicationUTF8)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("NO NEW COMMAND"))
				return
			}
			// Logging command
			entry := commands.CommandLog{
				CommandID: cmd.ID,
				Target:    t,
				Status:    "STZ_NEW",
				Action:    cmd.Action,
				Payload:   cmd.Payload,
			}
			if err := stzCommands.Log(&entry); err != nil {
				log.Printf("error logging new command %v", err)
			}
		}
	}
	// Send response
	httpResponse(w, http.StatusOK, CommandResponse{Message: "OK"})
}

// Handler to create new agent actions in the C2
func agentsActionHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	// Make sure the action is valid
	validActions := map[string]bool{
		"AGENT_DELETE": true,
	}
	action := r.PathValue("action")
	if action == "" {
		w.Header().Set(ContentType, JSONApplicationUTF8)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("NO ACTION"))
		return
	}
	if !validActions[strings.ToUpper(action)] {
		w.Header().Set(ContentType, JSONApplicationUTF8)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD ACTION"))
		return
	}
	// Read and decode POST body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var newAct StzAgentAction
	err = json.Unmarshal(body, &newAct)
	if err != nil {
		log.Println(err)
	}
	// Execute action
	switch a := newAct.Action; a {
	case "AGENT_DELETE":
		for _, a := range newAct.Agents {
			if err := stzAgents.Delete(a); err != nil {
				log.Println(err)
			}
		}
	case "AGENT_HIDE":
		for _, a := range newAct.Agents {
			if err := stzAgents.Hide(a); err != nil {
				log.Println(err)
			}
		}
	}
	// Send response
	httpResponse(w, http.StatusOK, CommandResponse{Message: "OK"})
}

// Handler to get all agents in JSON format
func jsonAgentsHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	// Retrieve agents
	agents, err := stzAgents.GetAllActive()
	if err != nil {
		log.Printf("error getting agents %v", err)
	}
	// Prepare data to be returned
	aJSON := []AgentJSON{}
	for _, a := range agents {
		nj := AgentJSON{
			ID:       a.ID,
			UUID:     a.UUID,
			Username: a.Username,
			Hostname: a.Hostname,
			IP:       a.IPsrc,
			GOOS:     a.GOOS,
			GOARCH:   a.GOARCH,
			TimeAgo:  pastTimeAgo(a.UpdatedAt),
		}
		aJSON = append(aJSON, nj)
	}
	returned := ReturnedAgents{
		Data: aJSON,
	}
	// Send response
	httpResponse(w, http.StatusOK, returned)
}

// Handler to do actions for the callbacks in the C2
func callbacksActionHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))
	// Make sure the action is valid
	validActions := map[string]bool{
		"CALLBACK_DELETE": true,
	}
	action := r.PathValue("action")
	if action == "" {
		w.Header().Set(ContentType, JSONApplicationUTF8)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("NO ACTION"))
		return
	}
	if !validActions[strings.ToUpper(action)] {
		w.Header().Set(ContentType, JSONApplicationUTF8)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD ACTION"))
		return
	}
	// Read and decode POST body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var newAct StzCallbackAction
	err = json.Unmarshal(body, &newAct)
	if err != nil {
		log.Println(err)
	}
	// Execute action
	switch a := newAct.Action; a {
	case "CALLBACK_DELETE":
		for _, a := range newAct.Callbacks {
			if err := stzCallbacks.Delete(a); err != nil {
				log.Println(err)
			}
		}
	case "CALLBACK_DISABLE":
		for _, a := range newAct.Callbacks {
			if err := stzCallbacks.Hide(a); err != nil {
				log.Println(err)
			}
		}
	}
	// Send response
	httpResponse(w, http.StatusOK, CommandResponse{Message: "OK"})
}
