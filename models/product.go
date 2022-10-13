package models

import "time"

type Product struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"product_name"`
	Category  string    `json:"category"`
	ImageURL  string    `json:"image"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
