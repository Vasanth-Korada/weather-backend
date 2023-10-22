package storage

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	if config.Host == "" || config.Port == "" || config.User == "" || config.Password == "" || config.DBName == "" || config.SSLMode == "" {
		return nil, errors.New("incomplete database configuration provided")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error creating a new database connection: %v", err)
		return nil, err
	}

	return db, nil
}
