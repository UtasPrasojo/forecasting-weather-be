// internal/models/activity.go
package models

import "time"

type Activity struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Name         string    `json:"name" example:"Rapat Koordinasi"`
    AreaCode     string    `json:"area_code" example:"31.71.01.1001"`
    
    Wilayah      Wilayah   `gorm:"foreignKey:AreaCode;references:Code" json:"wilayah"`

    ActivityDate time.Time `json:"activity_date" example:"2026-02-27T10:00:00Z"`

    WeatherStatus string    `json:"weather_status" swaggerignore:"true"`
    CreatedAt     time.Time `json:"created_at" swaggerignore:"true"`
    UpdatedAt     time.Time `json:"updated_at" swaggerignore:"true"`
}