package models

import "gorm.io/gorm"

type Weather struct {
	ID          uint     `gorm:"primary key;autoIncrement" json:"id"`
	Temperature *float32 `json:"temp"`
	City        *string  `json:"city"`
	Lat         *float32 `json:"lat"`
	Lng         *float32 `json:"lng"`
}

func MigrateWeather(db *gorm.DB) error {
	err := db.AutoMigrate(&Weather{})
	return err
}
