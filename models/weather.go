package models

import "gorm.io/gorm"

type Weather struct {
	ID          uint     `gorm:"primary key;autoIncrement" json:"id"`
	Temperature *float32 `json:"temp"`
	City        *string  `gorm:"unique;not null" json:"city"`
	Lat         *float32 `json:"lat"`
	Long        *float32 `json:"lng"`
	UserID      int64    `json:"user_id" gorm:"type:BIGINT"`
}

func MigrateWeather(db *gorm.DB) error {
	err := db.AutoMigrate(&Weather{})
	return err
}
