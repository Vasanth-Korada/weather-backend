package models

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	UserID    int    `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Password  string
	DOB       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func MigrateUsers(db *gorm.DB) error {
	if db == nil {
		return errors.New("database connection cannot be nil")
	}

	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Printf("Error migrating user model: %v", err)
		return err
	}

	return nil
}

func RegisterUser(db *gorm.DB, username, password, dob string) error {
	if db == nil {
		return errors.New("database connection cannot be nil")
	}
	if username == "" || password == "" || dob == "" {
		return errors.New("username, password, and dob fields cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	user := User{
		Username: username,
		Password: string(hashedPassword),
		DOB:      dob,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Printf("Error registering user: %v", err)
		return err
	}

	return nil
}

func AuthenticateUser(db *gorm.DB, username, password string) (*User, error) {
	if db == nil {
		return nil, errors.New("database connection cannot be nil")
	}
	if username == "" || password == "" {
		return nil, errors.New("username and password fields cannot be empty")
	}

	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		log.Printf("Error retrieving user: %v", err)
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		return nil, errors.New("incorrect password")
	}

	return &user, nil
}
