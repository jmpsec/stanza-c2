package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// DBString to format connection string to database for postgres
	DBString = "host=%s port=%s dbname=%s user=%s password=%s sslmode=disable"
	// DBKey to identify the configuration JSON key
	DBKey = "db"
)

// Get PostgreSQL DB using GORM
func getDB(cfg JSONConfigurationDB) (*gorm.DB, error) {
	postgresDSN := fmt.Sprintf(
		DBString, cfg.Host, cfg.Port, cfg.Name, cfg.Username, cfg.Password)
	dbConn, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := dbConn.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	// Performance settings for DB access
	// TODO: Make these configurable
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(30))
	return dbConn, nil
}
