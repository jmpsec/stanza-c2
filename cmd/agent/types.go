package main

import (
	"time"

	"github.com/jmpsec/stanza-c2/pkg/types"
)

// StzConfig to gather all the initialization parameters
type StzConfig struct {
	UUID         string
	Callbacks    []types.StzCallback
	Hostname     string
	LastCallback time.Time
	Registered   bool
}

// Worker to define struct of process
type Worker struct {
	Command string
	Output  chan string
}
