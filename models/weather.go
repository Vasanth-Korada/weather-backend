package models

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type Weather struct {
	ID          uint     `gorm:"primary key;autoIncrement" json:"id"`
	Temperature *float32 `json:"temp"`
	City        *string  `gorm:"unique;not null" json:"city"`
	Lat         *float32 `json:"lat"`
	Long        *float32 `json:"lng"`
	UserID      int64    `json:"user_id" gorm:"type:BIGINT"`
}

func MigrateWeather(db *gorm.DB) error {
	if db == nil {
		return errors.New("database connection cannot be nil")
	}

	if err := db.AutoMigrate(&Weather{}); err != nil {
		log.Printf("Error migrating Weather model: %v", err)
		return err
	}

	return nil
}
