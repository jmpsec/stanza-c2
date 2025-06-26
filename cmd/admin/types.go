package main

import (
	"github.com/jmpsec/stanza-c2/pkg/agents"
	"github.com/jmpsec/stanza-c2/pkg/commands"
	"github.com/jmpsec/stanza-c2/pkg/files"
	"github.com/jmpsec/stanza-c2/pkg/types"
)

// JSONConfigurationDB to hold all backend configuration values
type JSONConfigurationDB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// JSONConfigurationAdmin to hold all Admin server configuration values
type JSONConfigurationAdmin struct {
	Listener string `json:"listener"`
	Port     string `json:"port"`
	Host     string `json:"host"`
}

// StzNewCommand to create a new command
type StzNewCommand struct {
	Targets []string `json:"targets"`
	Action  string   `json:"action"`
	Payload string   `json:"payload"`
}

// StzAgentAction to create a new agent action
type StzAgentAction struct {
	Agents []string `json:"agents"`
	Action string   `json:"action"`
}

// StzCallbackAction to create a callback action
type StzCallbackAction struct {
	Callbacks []uint `json:"callbacks"`
	Action    string `json:"action"`
}

// ReturnedAgents to return a JSON with agents
type ReturnedAgents struct {
	Data []AgentJSON `json:"data"`
}

// AgentJSON to be used to populate JSON data for an agent
type AgentJSON struct {
	ID       uint   `json:"id"`
	Checkbox string `json:"checkbox"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Hostname string `json:"hostname"`
	IP       string `json:"ipsrc"`
	GOOS     string `json:"goos"`
	GOARCH   string `json:"goarch"`
	TimeAgo  string `json:"timeago"`
}

// CommandResponse to be returned to command requests
type CommandResponse struct {
	Message string `json:"message"`
}

// TableTemplateData to be used with the view of all agents
type TableTemplateData struct {
	Title                  string
	AgentsActiveClass      string
	AgentLogsActiveClass   string
	CallbacksActiveClass   string
	CommandsActiveClass    string
	CommandLogsActiveClass string
	FilesActiveClass       string
}

// AgentTemplateData to be used with the template for an agent
type AgentTemplateData struct {
	Title                  string
	AgentsActiveClass      string
	AgentLogsActiveClass   string
	CallbacksActiveClass   string
	CommandsActiveClass    string
	CommandLogsActiveClass string
	FilesActiveClass       string
	Details                agents.Agent
	Commands               []commands.Command
	Logs                   []commands.CommandLog
	Files                  []files.ExtractedFile
}

// AgentLogsTemplateData to be used with the template for agent logs
type AgentLogsTemplateData struct {
	Title                  string
	AgentsActiveClass      string
	AgentLogsActiveClass   string
	CallbacksActiveClass   string
	CommandsActiveClass    string
	CommandLogsActiveClass string
	FilesActiveClass       string
}

// CallbacksTemplateData to be used with the template for callbacks
type CallbacksTemplateData struct {
	Title                  string
	AgentsActiveClass      string
	AgentLogsActiveClass   string
	CallbacksActiveClass   string
	CommandsActiveClass    string
	CommandLogsActiveClass string
	FilesActiveClass       string
	Callbacks              []types.StzCallback
}

// CommandsTemplateData to be used with the template for callbacks
type CommandsTemplateData struct {
	Title                  string
	AgentsActiveClass      string
	AgentLogsActiveClass   string
	CallbacksActiveClass   string
	CommandsActiveClass    string
	CommandLogsActiveClass string
	FilesActiveClass       string
	Commands               []commands.Command
}

// CommandLogsTemplateData to be used with the template for callbacks
type CommandLogsTemplateData struct {
	Title                  string
	AgentsActiveClass      string
	AgentLogsActiveClass   string
	CallbacksActiveClass   string
	CommandsActiveClass    string
	CommandLogsActiveClass string
	FilesActiveClass       string
	Logs                   []commands.CommandLog
}

// OneLinersTemplateData to be used with the template for one-liners
type OneLinersTemplateData struct {
	Title                  string
	AgentsActiveClass      string
	AgentLogsActiveClass   string
	CallbacksActiveClass   string
	CommandsActiveClass    string
	CommandLogsActiveClass string
	FilesActiveClass       string
	OneLiners              []string
}
