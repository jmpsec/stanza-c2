package commands

import (
	"fmt"
	"log"

	"github.com/jmpsec/stanza-c2/pkg/types"
	"gorm.io/gorm"
)

// Command to keep all the commands to be sent to agents in beacons
type Command struct {
	gorm.Model
	Target    string `gorm:"index"`
	Action    string
	Payload   string
	Completed bool
}

// CommandLog to keep all the commands logs
type CommandLog struct {
	gorm.Model
	CommandID uint   `gorm:"index"`
	Target    string `gorm:"index"`
	Status    string
	Action    string
	Payload   string
	Data      string
}

// CommandManager to handle all commands
type CommandManager struct {
	DB *gorm.DB
}

// CreateCommandsManager to initialize the commands struct and its tables
func CreateCommandManager(backend *gorm.DB) *CommandManager {
	var c *CommandManager
	c = &CommandManager{DB: backend}
	if err := backend.AutoMigrate(Command{}); err != nil {
		log.Fatalf("Failed to AutoMigrate table (commands): %v", err)
	}
	if err := backend.AutoMigrate(CommandLog{}); err != nil {
		log.Fatalf("Failed to AutoMigrate table (command_logs): %v", err)
	}
	return c
}

// New to create a new command and assign it to its targets
func (c *CommandManager) New(command *Command) error {
	if err := c.DB.Create(&command).Error; err != nil {
		return fmt.Errorf("New command: %v", err)
	}
	return nil
}

// Log to add an entry to the log of commands
func (c *CommandManager) Log(entry *CommandLog) error {
	if err := c.DB.Create(&entry).Error; err != nil {
		return fmt.Errorf("Command log: %v", err)
	}
	return nil
}

// BeaconCommands to retrieve all new commands for an agent as response for a beacon
func (c *CommandManager) BeaconCommands(uuid string) ([]types.StzBeaconResponse, error) {
	result := []types.StzBeaconResponse{}
	commands := []Command{}
	if err := c.DB.Where("target = ? AND completed = ?", uuid, false).Find(&commands).Error; err != nil {
		return result, err
	}
	for _, c := range commands {
		result = append(result, types.StzBeaconResponse{
			ID:      c.ID,
			Action:  c.Action,
			Payload: c.Payload,
		})
	}
	return result, nil
}

// Get to retrieve a command by its ID
func (c *CommandManager) Get(cmdid uint) (Command, error) {
	var command Command
	if err := c.DB.Where("id = ?", cmdid).First(&command).Error; err != nil {
		return command, err
	}
	return command, nil
}

// GetAll to retrieve all the commands for an agent
func (c *CommandManager) GetAll(uuid string) ([]Command, error) {
	commands := []Command{}
	if err := c.DB.Order("updated_at desc").Where("target = ?", uuid).Find(&commands).Error; err != nil {
		return commands, err
	}
	return commands, nil
}

// GetLogs to retrieve all the commands for an agent
func (c *CommandManager) GetLogs(uuid string) ([]CommandLog, error) {
	logs := []CommandLog{}
	if err := c.DB.Order("created_at desc").Where("target = ?", uuid).Find(&logs).Error; err != nil {
		return logs, err
	}
	return logs, nil
}

// Update to update command and to add a CommandLog line to the DB
func (c *CommandManager) Update(cmdid uint, status, data string) error {
	// Set command as complete
	command, err := c.Get(cmdid)
	if err != nil {
		return fmt.Errorf("dbGetCommand %v", err)
	}
	if status == types.StzStatusDone || command.Action == "STZ_REGISTER" {
		if err := c.DB.Model(&command).Updates(map[string]interface{}{"completed": true}).Error; err != nil {
			return fmt.Errorf("Update %v", err)
		}
	}
	// Create entry in the logging table for commands
	entry := CommandLog{
		CommandID: command.ID,
		Target:    command.Target,
		Status:    status,
		Action:    command.Action,
		Payload:   command.Payload,
		Data:      data,
	}
	return c.Log(&entry)
}
