package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	// Header name to authenticate with Traphouse
	traphouseHeader = "X-Traphouse-Auth"
	// Token to authenticate with Traphouse
	traphouseToken = "this_should_be_the_token"
	// URL for Trapouse implants API
	traphouseAPIURL = "https://traphouse-domain/api/v1/implants/heartbeat"
)

// TraphouseImplant to score from C2
type TraphouseImplant struct {
	XID      string `json:"xid"`
	EgressIP string `json:"egress_ip"`
	OS       string `json:"os"`
	Hostname string `json:"hostname"`
	Username string `json:"user_info"`
}

// Report C2 implants to Traphouse to score more points
func traphouseImplantsScore(data TraphouseImplant) error {
	headers := map[string]string{
		traphouseHeader: traphouseToken,
	}
	jsonOut, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsonParam := strings.NewReader(string(jsonOut))
	log.Printf("-> TraphouseImplant : %v\n", jsonParam)
	resp, body, err := sendHTTPRequest("POST", traphouseAPIURL, jsonParam, headers)
	if resp != http.StatusOK {
		return errors.New("HTTP " + strconv.Itoa(resp) + ":" + string(body))
	}
	if err != nil {
		return err
	}
	return nil
}

/*
	// Send implant to Traphouse API if enabled
	if traphouseImplants {
		hasher := md5.New()
		hasher.Write([]byte(ip))
		XID := hex.EncodeToString(hasher.Sum(nil))
		implant := TraphouseImplant{
			XID:      XID,
			EgressIP: ip,
			OS:       "",
			Hostname: "",
			Username: "",
		}
		err = traphouseImplantsScore(implant)
		if err != nil {
			log.Println(err)
		}
		// Logging implant submission
		err = dbLoggingImplant(implant)
		if err != nil {
			log.Println(err)
		}
	}
*/
