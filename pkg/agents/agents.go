package agents

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Agent to keep the list of all registered agents
type Agent struct {
	gorm.Model
	UUID      string `gorm:"index"`
	LocalUUID string
	Username  string
	Hostname  string
	GOOS      string
	GOARCH    string
	Uname     string
	IPs       string
	IPsrc     string
	CycleMin  int
	CycleMax  int
	CycleNow  int
	Active    bool
}

// AgentLog to keep all the activity for each agent
type AgentLog struct {
	gorm.Model
	UUID   string `gorm:"index"`
	IPsrc  string
	Action string
	Data   string
}

// AgentManager to handle all agents
type AgentManager struct {
	DB *gorm.DB
}

// CreateAgentManager to initialize the agents struct and its tables
func CreateAgentManager(backend *gorm.DB) *AgentManager {
	var a *AgentManager
	a = &AgentManager{DB: backend}
	if err := backend.AutoMigrate(Agent{}); err != nil {
		log.Fatalf("Failed to AutoMigrate table (agents): %v", err)
	}
	if err := backend.AutoMigrate(AgentLog{}); err != nil {
		log.Fatalf("Failed to AutoMigrate table (agent_logs): %v", err)
	}
	return a
}

// Register to register an agent in the DB
func (a *AgentManager) Register(agent *Agent) error {
	if err := a.DB.Create(&agent).Error; err != nil {
		return fmt.Errorf("Register new agent: %v", err)
	}
	return nil
}

// Get to retrieve an agent from DB, by UUID
func (a *AgentManager) Get(uuid string) (Agent, error) {
	var agent Agent
	if err := a.DB.Where("uuid = ?", uuid).First(&agent).Error; err != nil {
		return agent, err
	}
	return agent, nil
}

// GetAllActive to retrieve all active agents
func (a *AgentManager) GetAllActive() ([]Agent, error) {
	var agents []Agent
	if err := a.DB.Where("active = ?", true).Find(&agents).Error; err != nil {
		return agents, err
	}
	return agents, nil
}

// Exist to check if an agent exists already in the DB
func (a *AgentManager) Exist(uuid, hostname, ipsrc, username string) (Agent, bool) {
	var agent Agent
	if err := a.DB.Where("active = ? AND uuid = ? AND hostname = ? AND ipsrc = ? AND username = ?", true, uuid, hostname, ipsrc, username).First(&agent).Error; err != nil {
		return agent, false
	}
	return agent, true
}

// ExistBeacon to check if an agent exists already in the DB
func (a *AgentManager) ExistBeacon(uuid, ipsrc string) (Agent, bool) {
	var agent Agent
	if err := a.DB.Where("active = ? AND uuid = ? AND ipsrc = ?", true, uuid, ipsrc).First(&agent).Error; err != nil {
		return agent, false
	}
	return agent, true
}

// UpdateBeaconCycle to update the current cycle value on beacon
func (a *AgentManager) UpdateBeaconCycle(uuid string, cycle int) error {
	agent, err := a.Get(uuid)
	if err != nil {
		return fmt.Errorf("dbGetAgent %v", err)
	}
	if err := a.DB.Model(&agent).Update("cycle_now", cycle).Error; err != nil {
		return fmt.Errorf("Update %v", err)
	}
	return nil
}

// Update to update an existing agent in the DB
func (a *AgentManager) Update(agentNew Agent) error {
	agent, err := a.Get(agentNew.UUID)
	if err != nil {
		return fmt.Errorf("dbGetAgent %v", err)
	}
	if err := a.DB.Model(&agent).Updates(agentNew).Error; err != nil {
		return fmt.Errorf("Updates %v", err)
	}
	return nil
}

// Log to add an entry to the log of agents
func (a *AgentManager) Log(entry *AgentLog) error {
	if err := a.DB.Create(&entry).Error; err != nil {
		return fmt.Errorf("Agent log: %v", err)
	}
	return nil
}

// LogCheckin to create an entry for the agent logging table and update agent last seen
func (a *AgentManager) LogCheckin(uuid, ip, action, data string) error {
	// Update agent time
	agent, err := a.Get(uuid)
	if err != nil {
		return fmt.Errorf("dbGetAgent %v", err)
	}
	agent.IPsrc = ip
	if err := a.DB.Model(&agent).Update("updated_at", time.Now()).Error; err != nil {
		return fmt.Errorf("Update %v", err)
	}
	// Create entry in the logging table for agents
	entry := AgentLog{
		UUID:   uuid,
		IPsrc:  ip,
		Action: action,
		Data:   data,
	}
	return a.Log(&entry)
}

// Delete to remove an agent from the C2
func (a *AgentManager) Delete(uuid string) error {
	agent, err := a.Get(uuid)
	if err != nil {
		return fmt.Errorf("dbGetAgent %v", err)
	}
	if err := a.DB.Delete(&agent).Error; err != nil {
		return fmt.Errorf("Delete %v", err)
	}
	return nil
}

// Hide an agent from the C2
func (a *AgentManager) Hide(uuid string) error {
	agent, err := a.Get(uuid)
	if err != nil {
		return fmt.Errorf("dbGetAgent %v", err)
	}
	if err := a.DB.Model(&agent).Update("active", false).Error; err != nil {
		return fmt.Errorf("Update %v", err)
	}
	return nil
}
