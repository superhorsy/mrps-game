package model

import (
	"time"
)

type Transaction struct {
	ID        uint32    `gorm:"primary_key"`
	UserID    uint32    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Reason    string    `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}
