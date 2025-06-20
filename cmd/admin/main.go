package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/jmpsec/stanza-c2/pkg/agents"
	"github.com/jmpsec/stanza-c2/pkg/callbacks"
	"github.com/jmpsec/stanza-c2/pkg/commands"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Define endpoints
const (
	// Service configuration file
	defConfigurationFile = "config/admin.json"
	// DB configuration key
	dbConfigKey = "db"
	// Service configuration key
	adminConfigKey = "admin"
)

// Private IP range with access - RFC 1918
var allowedNetworks = [3]string{"10.0.0.0/16", "192.168.0.0/16", "172.16.0.0/12"}

// Global variables
var (
	adminConfig  JSONConfigurationAdmin
	dbConfig     JSONConfigurationDB
	configFile   *string
	db           *gorm.DB
	stzAgents    *agents.AgentManager
	stzCallbacks *callbacks.CallbackManager
	stzCommands  *commands.CommandManager
)

// Checker for access via network ranges
func adminAllowed(ipaddress string) bool {
	for _, r := range allowedNetworks {
		_, subnet, _ := net.ParseCIDR(r)
		ip := net.ParseIP(ipaddress)
		if !subnet.Contains(ip) {
			return false
		}
	}
	return true
}

// Function to load the configuration file and assign to variables
func loadConfiguration(cfgFile string) (JSONConfigurationAdmin, JSONConfigurationDB, error) {
	log.Printf("Loading %s", cfgFile)
	var adminCfg JSONConfigurationAdmin
	var dbCfg JSONConfigurationDB
	// Load file and read config
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return adminCfg, dbCfg, err
	}
	// HTTP Service values
	adminRaw := viper.Sub(adminConfigKey)
	if adminRaw == nil {
		return adminCfg, dbCfg, fmt.Errorf("No configuration found for %s", adminConfigKey)
	}
	if err := adminRaw.Unmarshal(&adminCfg); err != nil {
		return adminCfg, dbCfg, err
	}
	// Backend values
	dbRaw := viper.Sub(dbConfigKey)
	if dbRaw == nil {
		return adminCfg, dbCfg, fmt.Errorf("No configuration found for %s", dbConfigKey)
	}
	if err := dbRaw.Unmarshal(&dbCfg); err != nil {
		return adminCfg, dbCfg, err
	}
	// No errors!
	return adminCfg, dbCfg, nil
}

// Get flags or environment variables and initialize token
func init() {
	// Check if configuration flag is provided as flag
	configFile = flag.String("config", defConfigurationFile, "Service configuration file")
	// Parse all defined flags
	flag.Parse()
	// Load service configuration
	var err error
	adminConfig, dbConfig, err = loadConfiguration(*configFile)
	if err != nil {
		log.Fatalf("Error loading configuration from %s [%s]", *configFile, err)
	}
}

// Go go!
func main() {
	var err error
	// Logging format flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	// Database handler
	db, err = getDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Database %s: %v", dbConfig.Name, err)
	}

	// Initialize agents manager
	stzAgents = agents.CreateAgentManager(db)
	// Initialize callbacks manager
	stzCallbacks = callbacks.CreateCallbackManager(db)
	// Initialize commands manager
	stzCommands = commands.CreateCommandManager(db)

	// Create router
	router := http.NewServeMux()

	// main
	router.HandleFunc("GET /", agentsHandler)
	// agents
	router.HandleFunc("GET /agents", agentsHandler)
	// agent logs
	router.HandleFunc("GET /agent-logs", agentLogsHandler)
	// commands
	router.HandleFunc("GET /commands", commandsHandler)
	// command logs
	router.HandleFunc("GET /command-logs", commandLogsHandler)
	// callbacks
	router.HandleFunc("GET /callbacks", callbacksHandler)
	// callbacks actions
	router.HandleFunc("POST /callbacks/{action}", callbacksActionHandler)
	// json handler for agents
	router.HandleFunc("GET /json/agents", jsonAgentsHandler)
	// agent view
	router.HandleFunc("GET /agent/{uuid}", agentViewHandler)
	// new commands
	router.HandleFunc("POST /commands/{action}", commandsActionHandler)
	// agent actions
	router.HandleFunc("POST /agents/{action}", agentsActionHandler)
	// commands output
	router.HandleFunc("GET /log/{uuid}", commandsActionHandler)
	// files and oneliners
	router.HandleFunc("GET /files-oneliners", filesOnelinersHandler)
	// favicon
	router.HandleFunc("GET /favicon.ico", faviconHandler)

	// static
	router.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	// files
	router.Handle("GET /files/", http.StripPrefix("/files", http.FileServer(http.Dir("./files"))))

	http.Handle("/", router)

	// Launch HTTP server
	service := adminConfig.Listener + ":" + adminConfig.Port
	log.Printf("STZ Admin: HTTP on port %s", service)
	log.Fatal(http.ListenAndServe(service, router))
}
