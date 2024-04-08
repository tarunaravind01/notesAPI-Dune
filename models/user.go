package models

import "gorm.io/gorm"

// gorm model for user table(postgres)
type User struct {
	gorm.Model
	Username 		string `gorm:"unique"`
	Password 		string
  }