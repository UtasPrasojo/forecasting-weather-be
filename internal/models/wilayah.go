package models

import (
	"time"
)

type Wilayah struct {
	Code      string    `gorm:"primaryKey;type:varchar(20)" json:"code"`
	Loc       string    `gorm:"type:varchar(255);not null" json:"loc"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}