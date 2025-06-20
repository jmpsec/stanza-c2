package main

// JSONConfigurationDB to hold all backend configuration values
type JSONConfigurationDB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// JSONConfigurationHTTP to hold all HTTP server configuration values
type JSONConfigurationHTTP struct {
	Listener      string `json:"listener"`
	Port          string `json:"port"`
	Host          string `json:"host"`
	CallbacksPort string `json:"callbacksPort"`
	Token         string `json:"token"`
}
