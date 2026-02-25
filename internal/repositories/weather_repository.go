package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"project-telkom-sigma/internal/models"
)

type WeatherRepository interface {
	UpsertWeather(weather *models.Weather) error
}

type weatherRepository struct {
	db *gorm.DB
}

func NewWeatherRepository(db *gorm.DB) WeatherRepository {
	return &weatherRepository{db: db}
}

func (r *weatherRepository) UpsertWeather(weather *models.Weather) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "area_code"}, {Name: "utc_datetime"}},
		DoUpdates: clause.AssignmentColumns([]string{"t", "hu", "weather_desc", "category", "sync_time"}),
	}).Create(weather).Error
}
