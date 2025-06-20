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
	// Default UDP port
	_udpPort = "6666"
	// Env variable for UDP port
	_udpPortEnv = "STZ_UDP_PORT"
	// Passcode to request callbacks
	_callbacksCode = "ThisShouldBeSecret"
	// Env variable for code to request callback
	_callbacksCodeEnv = "STZ_CALLBACKS_CODE"
	// Service configuration file
	configurationFile = "config/udp.json"
)

// Global variables
var (
	token   string
	udpPort = _udpPort
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
	udpPort = getServerEnv(_udpPortEnv, _udpPort)
}

// Handles incoming connections.
func handleConnection(conn *net.UDPConn) {
	var buf [1024]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	conn.WriteToUDP([]byte("STZ: Message received."), addr)
}

// Go go!
func main() {
	log.Printf("Single token for this sessions is %s\n", token)

	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+udpPort)
	if err != nil {
		log.Fatal("Error resolving", err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Error listening", err)
	}

	// Launch UDP server
	log.Printf("STZ UDP Server should be listening on port %s", udpPort)
	for {
		handleConnection(conn)
	}
}
