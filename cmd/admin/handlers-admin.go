package main

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
)

// Handle error requests
func errorHTTPHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	t, _ := template.ParseFiles("templates/error.html")
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("template error %v", err)
		return
	}
}

// Handler to display the main table with all agents
func agentsHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if !adminAllowed(ip) {
		log.Printf("ACL: %s has no access to admin", ip)
		//return
	}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	t, err := template.ParseFiles(
		"templates/table.html",
		"templates/page-js.html",
		"templates/page-footer.html",
		"templates/page-modals.html",
		"templates/page-sidebar.html",
		"templates/page-header.html",
		"templates/page-head.html",
	)
	if err != nil {
		log.Printf("error getting agents template: %v", err)
		return
	}
	templateData := TableTemplateData{
		Title:             "Stanza C2: Agents",
		AgentsActiveClass: "active",
	}
	if err := t.Execute(w, templateData); err != nil {
		log.Printf("template error %v", err)
		return
	}
}

// Handler to display all the agent logs
func agentLogsHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if !adminAllowed(ip) {
		log.Printf("ACL: %s has no access to admin", ip)
		//return
	}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	t, err := template.ParseFiles(
		"templates/agent-logs.html",
		"templates/page-js.html",
		"templates/page-footer.html",
		"templates/page-modals.html",
		"templates/page-sidebar.html",
		"templates/page-header.html",
		"templates/page-head.html",
	)
	if err != nil {
		log.Printf("error getting agents template: %v", err)
		return
	}
	templateData := AgentLogsTemplateData{
		Title:                "Stanza C2: Agents Logs",
		AgentLogsActiveClass: "active",
	}
	if err := t.Execute(w, templateData); err != nil {
		log.Printf("template error %v", err)
		return
	}
}

// Handler to display all the callbacks
func callbacksHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if !adminAllowed(ip) {
		log.Printf("ACL: %s has no access to admin", ip)
		//return
	}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))
	calls, err := stzCallbacks.StzCallbacks()
	if err != nil {
		log.Println(err)
	}
	t, err := template.ParseFiles(
		"templates/callbacks.html",
		"templates/page-js.html",
		"templates/page-footer.html",
		"templates/page-modals.html",
		"templates/page-sidebar.html",
		"templates/page-header.html",
		"templates/page-head.html",
	)
	if err != nil {
		log.Printf("error getting agents template: %v", err)
		return
	}
	templateData := CallbacksTemplateData{
		Title:                "Stanza C2: Callbacks",
		CallbacksActiveClass: "active",
		Callbacks:            calls,
	}
	if err := t.Execute(w, templateData); err != nil {
		log.Printf("template error %v", err)
		return
	}
}

// Handler to display all the commands
func commandsHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if !adminAllowed(ip) {
		log.Printf("ACL: %s has no access to admin", ip)
		//return
	}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	t, err := template.ParseFiles(
		"templates/commands.html",
		"templates/page-js.html",
		"templates/page-footer.html",
		"templates/page-modals.html",
		"templates/page-sidebar.html",
		"templates/page-header.html",
		"templates/page-head.html",
	)
	if err != nil {
		log.Printf("error getting agents template: %v", err)
		return
	}
	templateData := CommandsTemplateData{
		Title:               "Stanza C2: Commands",
		CommandsActiveClass: "active",
	}
	if err := t.Execute(w, templateData); err != nil {
		log.Printf("template error %v", err)
		return
	}
}

// Handler to display all the command logs
func commandLogsHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if !adminAllowed(ip) {
		log.Printf("ACL: %s has no access to admin", ip)
		//return
	}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	t, err := template.ParseFiles(
		"templates/command-logs.html",
		"templates/page-js.html",
		"templates/page-footer.html",
		"templates/page-modals.html",
		"templates/page-sidebar.html",
		"templates/page-header.html",
		"templates/page-head.html",
	)
	if err != nil {
		log.Printf("error getting agents template: %v", err)
		return
	}
	templateData := CommandLogsTemplateData{
		Title:                  "Stanza C2: Command Logs",
		CommandLogsActiveClass: "active",
	}
	if err := t.Execute(w, templateData); err != nil {
		log.Printf("template error %v", err)
		return
	}
}

// Handler to display a single agent view with all details and actions
func agentViewHandler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n", string(requestDump))

	uuid := r.PathValue("uuid")
	if uuid == "" {
		return
	}
	// Retrieve agent
	agent, err := stzAgents.Get(uuid)
	if err != nil {
		log.Printf("error getting agent %v", err)
		return
	}
	// Retrieve logs
	logs, err := stzCommands.GetLogs(uuid)
	if err != nil {
		log.Printf("error getting logs %v", err)
		return
	}
	// Retrieve commands
	commands, err := stzCommands.GetAll(uuid)
	if err != nil {
		log.Printf("error getting commands %v", err)
		return
	}
	// Retrieve files
	files, err := stzFiles.GetAll(uuid)
	// Custom functions to handle formatting
	funcMap := template.FuncMap{
		"pastTimeAgo":   pastTimeAgo,
		"payloadFormat": payloadFormat,
	}
	// Prepare template
	t, err := template.New("agent.html").Funcs(funcMap).ParseFiles(
		"templates/agent.html",
		"templates/page-js.html",
		"templates/page-footer.html",
		"templates/page-modals.html",
		"templates/page-sidebar.html",
		"templates/page-header.html",
		"templates/page-head.html",
	)
	if err != nil {
		log.Printf("error getting agents template: %v", err)
		return
	}
	data := AgentTemplateData{
		Title:    "Stanza C2: Agent view",
		Details:  agent,
		Commands: commands,
		Logs:     logs,
		Files:    files,
	}
	// Execute template
	if err := t.Execute(w, data); err != nil {
		log.Printf("template error %v", err)
		return
	}
}

// Handler to display files and oneliners
func filesOnelinersHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if !adminAllowed(ip) {
		log.Printf("ACL: %s has no access to admin", ip)
		//return
	}
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s\n\n", string(requestDump))

	t, err := template.ParseFiles(
		"templates/files-oneliners.html",
		"templates/page-js.html",
		"templates/page-footer.html",
		"templates/page-modals.html",
		"templates/page-sidebar.html",
		"templates/page-header.html",
		"templates/page-head.html",
	)
	if err != nil {
		log.Printf("error getting agents template: %v", err)
		return
	}
	templateData := OneLinersTemplateData{
		Title:            "Stanza C2: Files and Oneliners",
		FilesActiveClass: "active",
		OneLiners:        []string{"stanza-admin-dev"},
	}
	if err := t.Execute(w, templateData); err != nil {
		log.Printf("template error %v", err)
		return
	}
}
