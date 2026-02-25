package models

import "time"

type Activity struct {
	ID            uint      `gorm:"primaryKey"`
	Name          string    `json:"name"`           // Nama kegiatan (misal: "Meeting Telkom")
	AreaCode      string    `json:"area_code"`      // Relasi ke Wilayah.Code (adm4)
	ActivityDate  time.Time `json:"activity_date"`  // Tanggal dan Jam kegiatan
	WeatherStatus string    `json:"weather_status"` // Akan diisi otomatis oleh sistem
	SyncTime      time.Time `json:"sync_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
