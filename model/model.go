package model

import (
	"time"
)

// Tax model
type Tax struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	TaxCode   uint      `json:"tax_code"`
}

// TaxCode model for Tax
type TaxCode struct {
	Code uint   `json:"code" gorm:"primary_key"`
	Name string `json:"name"`
}
