package services

import (
	"errors"
	"project-telkom-sigma/internal/models"
	"project-telkom-sigma/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type ActivityService interface {
	CreateActivity(input *models.Activity) (*models.Activity, error)
	GetAllActivities(search, sortBy, order string) ([]models.Activity, error)
	UpdateActivity(id string, input *models.Activity) (*models.Activity, error)
	DeleteActivity(id string) error
}

type activityService struct {
	db          *gorm.DB
	weatherRepo repositories.WeatherRepository
}

func NewActivityService(db *gorm.DB, wRepo repositories.WeatherRepository) ActivityService {
	return &activityService{
		db:          db,
		weatherRepo: wRepo,
	}
}

func (s *activityService) CreateActivity(input *models.Activity) (*models.Activity, error) {
	input.WeatherStatus = s.getWeatherStatus(input.AreaCode, input.ActivityDate)

	if err := s.db.Create(input).Error; err != nil {
		return nil, err
	}
	return input, nil
}

func (s *activityService) UpdateActivity(id string, input *models.Activity) (*models.Activity, error) {
	var activity models.Activity
	if err := s.db.First(&activity, id).Error; err != nil {
		return nil, errors.New("kegiatan tidak ditemukan")
	}

	activity.Name = input.Name
	activity.AreaCode = input.AreaCode
	activity.ActivityDate = input.ActivityDate

	activity.WeatherStatus = s.getWeatherStatus(activity.AreaCode, activity.ActivityDate)

	if err := s.db.Save(&activity).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *activityService) getWeatherStatus(areaCode string, date time.Time) string {
	var weather models.Weather

	err := s.db.Where("area_code = ? AND local_datetime <= ?", areaCode, date).
		Order("local_datetime DESC").
		First(&weather).Error

	if err == nil {
		return weather.WeatherDesc
	}
	return "Cuaca tidak diketahui (Silakan Sync data BMKG)"
}


func (s *activityService) GetAllActivities(search, sortBy, order string) ([]models.Activity, error) {
    var activities []models.Activity
    
    dbQuery := s.db.Model(&models.Activity{}).Preload("Wilayah")

    if search != "" {
        searchText := "%" + search + "%"
        dbQuery = dbQuery.Where("name ILIKE ? OR weather_status ILIKE ?", searchText, searchText)
    }

    if err := dbQuery.Order(sortBy + " " + order).Find(&activities).Error; err != nil {
        return nil, err
    }
    return activities, nil
}

func (s *activityService) DeleteActivity(id string) error {
	return s.db.Delete(&models.Activity{}, id).Error
}
