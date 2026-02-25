package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-telkom-sigma/internal/models"
	"project-telkom-sigma/internal/repositories"
	"time"
)

type WeatherService interface {
	SyncWeather(adm4 string) error
}

type weatherService struct {
	repo repositories.WeatherRepository
}

func NewWeatherService(r repositories.WeatherRepository) WeatherService {
	return &weatherService{repo: r}
}

func (s *weatherService) SyncWeather(adm4 string) error {
	url := fmt.Sprintf("https://api.bmkg.go.id/publik/prakiraan-cuaca?adm4=%s", adm4)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var bmkg models.BMKGResponse
	if err := json.NewDecoder(resp.Body).Decode(&bmkg); err != nil {
		return err
	}

	for _, dataItem := range bmkg.Data {
		for _, cuacaArray := range dataItem.Cuaca {
			for _, item := range cuacaArray {
				utcTime, _ := time.Parse("2006-01-02 15:04:05", item.UtcDatetime)
				localTime, _ := time.Parse("2006-01-02 15:04:05", item.LocalDatetime)

				// Logika Bisnis Penentuan Kategori
				category := "Lainnya"
				desc := item.WeatherDesc
				if desc == "Cerah" || desc == "Cerah Berawan" {
					category = "Cerah"
				} else if desc == "Berawan" || desc == "Berawan Tebal" {
					category = "Berawan"
				} else if desc == "Hujan Ringan" || desc == "Hujan Sedang" || desc == "Hujan Petir" {
					category = "Hujan"
				}

				weatherEntry := models.Weather{
					AreaCode:      dataItem.Lokasi.Adm4,
					UtcDatetime:   utcTime,
					LocalDatetime: localTime,
					T:             item.T,
					Hu:            item.Hu,
					WeatherDesc:   item.WeatherDesc,
					Category:      category,
					SyncTime:      time.Now(),
				}

				// Simpan melalui Repository
				if err := s.repo.UpsertWeather(&weatherEntry); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
