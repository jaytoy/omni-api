package models

import "time"

type User struct {
	ID               uint   `json:"id" gorm:"primary_key"`
	FirstName        string `json:"first_name" gorm:"not null"`
	LastName         string `json:"last_name" gorm:"not null"`
	Email            string `json:"email" gorm:"uniqueIndex;not null"`
	Password         string
	VerificationCode string
	Verified         bool
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
