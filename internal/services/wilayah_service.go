package services

import (
	"project-telkom-sigma/internal/models"
	"project-telkom-sigma/internal/repositories"
)

type WilayahService interface {
	SearchWilayah(query string) ([]models.Wilayah, error)
}

type wilayahService struct {
	repo repositories.WilayahRepository
}

func NewWilayahService(r repositories.WilayahRepository) WilayahService {
	return &wilayahService{repo: r}
}

func (s *wilayahService) SearchWilayah(query string) ([]models.Wilayah, error) {
	return s.repo.GetAllWilayah(query)
}