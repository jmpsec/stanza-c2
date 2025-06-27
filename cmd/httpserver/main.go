package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jmpsec/stanza-c2/pkg/agents"
	"github.com/jmpsec/stanza-c2/pkg/callbacks"
	"github.com/jmpsec/stanza-c2/pkg/commands"
	"github.com/jmpsec/stanza-c2/pkg/files"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Define endpoints
const (
	// Redirect URL
	redirectURL = "https://www.nccdc.org"
	// Default endpoint to handle callbacks retrieval
	callbacksPath = "/__c"
	// Default endpoint to handle HTTP registrations
	registerPath = "/__r"
	// Default endpoint to handle HTTP beacons
	beaconPath = "/__b"
	// Default endpoint to handle HTTP execution
	executionPath = "/__x"
	// Default endpoint to receive files with HTTP
	filesPath = "/__f"
	// Service configuration file
	defConfigurationFile = "config/http.json"
	// DB configuration key
	dbConfigKey = "db"
	// Service configuration key
	httpConfigKey = "http"
	// Default files folder
	defFilesFolder = "./files"
)

// Global variables
var (
	httpConfig   JSONConfigurationHTTP
	dbConfig     JSONConfigurationDB
	configFile   *string
	db           *gorm.DB
	stzAgents    *agents.AgentManager
	stzCallbacks *callbacks.CallbackManager
	stzCommands  *commands.CommandManager
	stzFiles     *files.FileManager
)

// Function to load the configuration file and assign to variables
func loadConfiguration(cfgFile string) (JSONConfigurationHTTP, JSONConfigurationDB, error) {
	log.Printf("Loading %s", cfgFile)
	var httpCfg JSONConfigurationHTTP
	var dbCfg JSONConfigurationDB
	// Load file and read config
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return httpCfg, dbCfg, err
	}
	// HTTP Service values
	httpRaw := viper.Sub(httpConfigKey)
	if httpRaw == nil {
		return httpCfg, dbCfg, fmt.Errorf("No configuration found for %s", httpConfigKey)
	}
	if err := httpRaw.Unmarshal(&httpCfg); err != nil {
		return httpCfg, dbCfg, err
	}
	// Backend values
	dbRaw := viper.Sub(dbConfigKey)
	if dbRaw == nil {
		return httpCfg, dbCfg, fmt.Errorf("No configuration found for %s", dbConfigKey)
	}
	if err := dbRaw.Unmarshal(&dbCfg); err != nil {
		return httpCfg, dbCfg, err
	}
	// No errors!
	return httpCfg, dbCfg, nil
}

// Get flags or environment variables and initialize token
func init() {
	// Check if configuration flag is provided as flag
	configFile = flag.String("config", defConfigurationFile, "Service configuration file")
	// Parse all defined flags
	flag.Parse()
	// Load service configuration
	var err error
	httpConfig, dbConfig, err = loadConfiguration(*configFile)
	if err != nil {
		log.Fatalf("Error loading service configuration from %s [%s]", *configFile, err)
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
	// Initialize file manager
	stzFiles = files.CreateFileManager(db)

	// Check if callbacks/endpoints are set by the JSON host
	if !stzCallbacks.CheckByHost(httpConfig.Host) {
		log.Printf("No callbacks found for %s", httpConfig.Host)
		if err := stzCallbacks.New(httpConfig.Host, httpConfig.CallbacksPort, "http"); err != nil {
			log.Fatalf("Failed to initialize callback for HTTP %v", err)
		}
		//if err := dbNewCallback(httpConfig.Host, "443", "https://"); err != nil {
		//	log.Fatalf("Failed to initialize callback for HTTPS %v", err)
		//}
	}
	log.Printf("Callbacks ready for %s", httpConfig.Host)

	// Create router
	router := http.NewServeMux()

	// Handlers for different operations
	router.HandleFunc("/", emptyHTTPHandler)

	// register
	router.HandleFunc("GET "+registerPath, redirectHTTPHandler)
	router.HandleFunc("POST "+registerPath, registerHTTPHandler)

	// beacon
	router.HandleFunc("GET "+beaconPath, redirectHTTPHandler)
	router.HandleFunc("POST "+beaconPath, beaconHTTPHandler)

	// execution
	router.HandleFunc("GET "+executionPath, redirectHTTPHandler)
	router.HandleFunc("POST "+executionPath, executionHTTPHandler)

	// callbacks
	router.HandleFunc("GET "+callbacksPath, redirectHTTPHandler)
	router.HandleFunc("POST "+callbacksPath, callbacksHTTPHandler)

	// files
	router.Handle("GET /files/", http.StripPrefix("/files", http.FileServer(http.Dir(defFilesFolder))))
	router.HandleFunc("POST "+filesPath, filesHTTPHandler)

	http.Handle("/", router)

	// Launch HTTP server for localhost
	service := httpConfig.Listener + ":" + httpConfig.Port
	log.Printf("STZ: HTTP on port %s", service)
	log.Fatal(http.ListenAndServe(service, router))
}
