package models

import (
	"time"
)

type Wilayah struct {
	// Kita gunakan Code sebagai Primary Key karena nilainya unik (misal: 11.01.01)
	Code      string    `gorm:"primaryKey;type:varchar(20)" json:"code"`
	Loc       string    `gorm:"type:varchar(255);not null" json:"loc"`
	
	// Tambahan untuk audit trail (opsional tapi sangat disarankan untuk tugas)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}