package callbacks

import (
	"fmt"
	"log"

	"github.com/jmpsec/stanza-c2/pkg/types"
	"gorm.io/gorm"
)

const (
	registerEndpoint  = "register"
	beaconEndpoint    = "beacon"
	callbacksEndpoint = "callbacks"
	executionEndpoint = "execution"
	filesEndpoint     = "files"
)

const (
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
)

var _endpoints = map[string]string{
	registerEndpoint:  registerPath,
	beaconEndpoint:    beaconPath,
	callbacksEndpoint: callbacksPath,
	executionEndpoint: executionPath,
	filesEndpoint:     filesPath,
}

// Endpoint to keep the list of endpoints per callback
type Endpoint struct {
	gorm.Model
	CallbackID uint `gorm:"index"`
	Type       string
	Value      string
	Active     bool
}

// Callback to keep the list of callbacks to be used
type Callback struct {
	gorm.Model
	Host     string
	Port     string
	Protocol string
	Active   bool
}

// CallbackManager to handle all callbacks and endpoints
type CallbackManager struct {
	DB *gorm.DB
}

// CreateCallbackManager to initialize the callbacks struct and its tables
func CreateCallbackManager(backend *gorm.DB) *CallbackManager {
	var c *CallbackManager
	c = &CallbackManager{DB: backend}
	if err := backend.AutoMigrate(Callback{}); err != nil {
		log.Fatalf("Failed to AutoMigrate table (callbacks): %v", err)
	}
	if err := backend.AutoMigrate(Endpoint{}); err != nil {
		log.Fatalf("Failed to AutoMigrate table (endpoints): %v", err)
	}
	return c
}

// GetEndpoints to retrieve all the endpoints of a callback
func (c *CallbackManager) GetEndpoints(callbackid uint) (map[string]string, error) {
	var endpoints []Endpoint
	result := map[string]string{}
	if err := c.DB.Where("callback_id = ?", callbackid).Find(&endpoints).Error; err != nil {
		return result, err
	}
	for _, e := range endpoints {
		result[e.Type] = e.Value
	}
	return result, nil
}

// CheckByHost to check if callbacks are ready by host
func (c *CallbackManager) CheckByHost(host string) bool {
	var results int64
	c.DB.Model(&Callback{}).Where("host = ?", host).Count(&results)
	return (results > 0)
}

// GetAll to retrieve all the callbacks for a server
func (c *CallbackManager) GetAll() ([]Callback, error) {
	var callbacks []Callback
	if err := c.DB.Find(&callbacks).Error; err != nil {
		return callbacks, err
	}
	return callbacks, nil
}

// Helper to generate HTTP endpoints for a host
func (c *CallbackManager) generateEndpoints(host, port, protocol string) map[string]string {
	res := make(map[string]string)
	for k, v := range _endpoints {
		res[k] = protocol + "://" + host + ":" + port + v
	}
	return res
}

// CreateEndpoints to generate all the endpoints for one callback
func (c *CallbackManager) CreateEndpoints(callbackid uint, host, port, protocol string) error {
	endpoints := c.generateEndpoints(host, port, protocol)
	for eType, eValue := range endpoints {
		if err := c.NewEndpoint(callbackid, eType, eValue); err != nil {
			return fmt.Errorf("error creating endpoint for %s - %v", eType, err)
		}
	}
	return nil
}

// New to create a new callback with its endpoints
func (c *CallbackManager) New(host, port, protocol string) error {
	entry := Callback{
		Host:     host,
		Port:     port,
		Protocol: protocol,
		Active:   true,
	}
	if err := c.DB.Create(&entry).Error; err != nil {
		return fmt.Errorf("Create Callback %v", err)
	}
	if err := c.CreateEndpoints(entry.ID, host, port, protocol); err != nil {
		return err
	}
	return nil
}

// NewEndpoint to create a new endpoint to be used with a callback
func (c *CallbackManager) NewEndpoint(callbackid uint, eType, eValue string) error {
	entry := Endpoint{
		CallbackID: callbackid,
		Type:       eType,
		Value:      eValue,
		Active:     true,
	}
	if err := c.DB.Create(&entry).Error; err != nil {
		return fmt.Errorf("Create Endpoint %v", err)
	}
	return nil
}

// StzCallbacks to extract callbacks and endpoints and prepares the data to be used
func (c *CallbackManager) StzCallbacks() ([]types.StzCallback, error) {
	callbacks, err := c.GetAll()
	result := []types.StzCallback{}
	if err != nil {
		return result, fmt.Errorf("GetAll %v", err)

	}
	for _, _c := range callbacks {
		endp, err := c.GetEndpoints(_c.ID)
		if err != nil {
			return result, fmt.Errorf("GetEndpoints %v", err)
		}
		stz := types.StzCallback{
			ID:        _c.ID,
			Host:      _c.Host,
			Port:      _c.Port,
			Protocol:  _c.Protocol,
			Endpoints: endp,
		}
		result = append(result, stz)
	}
	return result, nil
}

// Delete to remove a callback from the C2
func (c *CallbackManager) Delete(callbackid uint) error {
	var endpoints []Endpoint
	if err := c.DB.Where("callback_id = ?", callbackid).Find(&endpoints).Error; err != nil {
		return err
	}
	for _, e := range endpoints {
		if err := c.DB.Delete(&e).Error; err != nil {
			return fmt.Errorf("Delete %v", err)
		}
	}
	var callback Callback
	if err := c.DB.Where("id = ?", callbackid).First(&callback).Error; err != nil {
		return err
	}
	if err := c.DB.Delete(&callback).Error; err != nil {
		return fmt.Errorf("Delete %v", err)
	}
	return nil
}

// Hide to disable a callback from the C2
func (c *CallbackManager) Hide(callbackid uint) error {
	return nil
}
