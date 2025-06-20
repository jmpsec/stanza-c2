package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

// Define endpoints
const (
	// Default TCP port
	_tcpPort = "5555"
	// Env variable for TCP port
	_tcpPortEnv = "STZ_TCP_PORT"
	// Passcode to request callbacks
	_callbacksCode = "ThisShouldBeSecret"
	// Env variable for code to request callback
	_callbacksCodeEnv = "STZ_CALLBACKS_CODE"
	// Service configuration file
	configurationFile = "config/tcp.json"
)

// Global variables
var (
	token   string
	tcpPort = _tcpPort
)

// Helper to get environment variables
func getServerEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Get flags or environment variables and initialize token
func init() {
	// Generate token
	hasher := md5.New()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	hasher.Write([]byte(strconv.Itoa(rand.Int())))
	token = hex.EncodeToString(hasher.Sum(nil))

	// Environment variables or flags
	tcpPort = getServerEnv(_tcpPortEnv, _tcpPort)
}

// Handles incoming connections.
func handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	// Close connection before exiting
	defer conn.Close()

	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading:", err.Error())
	}
	log.Println(conn.RemoteAddr())

	// Send a response back
	conn.Write([]byte("STZ: Message received."))
}

// Go go!
func main() {
	log.Printf("Single token for this session is %s\n", token)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+tcpPort)
	if err != nil {
		log.Fatal("Error resolving", err)
	}
	// Listen for incoming connections.
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	// Close the listener when the application closes.
	defer listener.Close()

	// Launch TCP server
	log.Printf("STZ TCP Server should be listening on port %s", tcpPort)
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting: ", err.Error())
		}
		// Handle connections in a new goroutine.
		go handleConnection(conn)
	}
}
