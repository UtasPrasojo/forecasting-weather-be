package repositories

import (
	"project-telkom-sigma/internal/models"
	"gorm.io/gorm"
)

type WilayahRepository interface {
	GetAllWilayah(search string) ([]models.Wilayah, error)
}

type wilayahRepository struct {
	db *gorm.DB
}

func NewWilayahRepository(db *gorm.DB) WilayahRepository {
	return &wilayahRepository{db: db}
}

func (r *wilayahRepository) GetAllWilayah(search string) ([]models.Wilayah, error) {
	var wilayahs []models.Wilayah
	query := r.db.Model(&models.Wilayah{})

	if search != "" {
		query = query.Where("code LIKE ? OR loc ILIKE ?", search+"%", "%"+search+"%")
	}

	err := query.Limit(20).Find(&wilayahs).Error 
	return wilayahs, err
}