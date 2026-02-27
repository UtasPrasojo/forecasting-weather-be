package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/models"
	"project-telkom-sigma/internal/services"
	"time"
)

type ActivityHandler struct {
	WeatherService services.WeatherService 
}

func NewActivityHandler(s services.WeatherService) *ActivityHandler {
	return &ActivityHandler{WeatherService: s}
}

// CreateActivity godoc
// @Summary      Buat Rencana Kegiatan
// @Description  Membuat rencana kegiatan baru dan otomatis mencocokkan status cuaca berdasarkan data BMKG yang tersedia
// @Tags         Activity
// @Accept       json
// @Produce      json
// @Param        activity  body      models.Activity  true  "Data JSON Kegiatan"
// @Success      201       {object}  models.Activity
// @Failure      400       {object}  map[string]string
// @Router       /api/activity [post]
func (h *ActivityHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var input models.Activity
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
		return
	}

	var lastWeather models.Weather
	if err := database.DB.Order("sync_time DESC").First(&lastWeather).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Belum ada data wilayah yang di-sync. Silakan lakukan Sync terlebih dahulu.",
		})
		return
	}

	input.AreaCode = lastWeather.AreaCode
	log.Printf("Menggunakan wilayah otomatis dari Sync terakhir: %s", input.AreaCode)

	var weather models.Weather
	findWeather := func(code string, date time.Time) error {
		return database.DB.Where("area_code = ?", code).
			Order(fmt.Sprintf("ABS(EXTRACT(EPOCH FROM (local_datetime - '%s')))",
				date.Format("2006-01-02 15:04:05"))).
			First(&weather).Error
	}

	err := findWeather(input.AreaCode, input.ActivityDate)

	if err == nil {
		input.WeatherStatus = weather.WeatherDesc
		log.Printf("Cuaca ditemukan: %s", weather.WeatherDesc)
	} else {
		log.Printf("Data cuaca jam tersebut tidak ada, mencoba Sync ulang area: %s", input.AreaCode)
		h.WeatherService.SyncWeather(input.AreaCode)

		if errRetry := findWeather(input.AreaCode, input.ActivityDate); errRetry == nil {
			input.WeatherStatus = weather.WeatherDesc
		} else {
			input.WeatherStatus = "Cuaca tidak diketahui (Gagal Sync Ulang)"
		}
	}

	if err := database.DB.Create(&input).Error; err != nil {
		log.Printf("Error DB: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menyimpan kegiatan"})
		return
	}

	writeJSON(w, http.StatusCreated, input)
}

// GetAllActivities godoc
// @Summary      Daftar Semua Kegiatan
// @Description  Mengambil semua daftar kegiatan dengan fitur pencarian dan pengurutan
// @Tags         Activity
// @Produce      json
// @Param        search   query     string  false  "Cari berdasarkan nama atau status cuaca"
// @Param        sort_by  query     string  false  "Kolom pengurutan (contoh: activity_date, name)"
// @Param        order    query     string  false  "Urutan data (asc/desc)"
// @Success      200      {array}   models.Activity
// @Failure      500      {object}  map[string]string
// @Router       /api/activity [get]
func (h *ActivityHandler) GetAllActivities(w http.ResponseWriter, r *http.Request) {
	var activities []models.Activity

	query := r.URL.Query()
	search := query.Get("search")
	sortBy := query.Get("sort_by")
	order := query.Get("order")

	dbQuery := database.DB.Model(&models.Activity{}).Preload("Wilayah")

	if search != "" {
		searchText := "%" + search + "%"
		dbQuery = dbQuery.Where("name ILIKE ? OR weather_status ILIKE ?", searchText, searchText)
	}

	if sortBy == "" {
		sortBy = "activity_date"
	}
	if order == "" || (order != "asc" && order != "desc") {
		order = "desc"
	}

	err := dbQuery.Order(sortBy + " " + order).Find(&activities).Error
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil data kegiatan"})
		return
	}

	writeJSON(w, http.StatusOK, activities)
}

// UpdateActivity godoc
// @Summary      Perbarui Kegiatan
// @Description  Mengubah data kegiatan dan memperbarui status cuaca secara otomatis
// @Tags         Activity
// @Accept       json
// @Produce      json
// @Param        id        query     string           true  "ID Kegiatan yang akan diubah"
// @Param        activity  body      models.Activity  true  "Data JSON Kegiatan baru"
// @Success      200       {object}  models.Activity
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Router       /api/activity/update [put]
func (h *ActivityHandler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "ID wajib diisi"})
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, id).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "Kegiatan tidak ditemukan"})
		return
	}

	var input models.Activity
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
		return
	}

	activity.Name = input.Name
	activity.AreaCode = input.AreaCode
	activity.ActivityDate = input.ActivityDate

	var weather models.Weather
	err := database.DB.Where("area_code = ?", activity.AreaCode).
		Order(fmt.Sprintf("ABS(EXTRACT(EPOCH FROM (local_datetime - '%s')))",
			activity.ActivityDate.Format("2006-01-02 15:04:05"))).
		First(&weather).Error

	if err == nil {
		activity.WeatherStatus = weather.WeatherDesc
	} else {
		log.Printf("Data tidak ada, mencoba sync ulang untuk area: %s", activity.AreaCode)
		h.WeatherService.SyncWeather(activity.AreaCode)

		database.DB.Where("area_code = ?", activity.AreaCode).
			Order(fmt.Sprintf("ABS(EXTRACT(EPOCH FROM (local_datetime - '%s')))",
				activity.ActivityDate.Format("2006-01-02 15:04:05"))).
			First(&weather)

		if weather.WeatherDesc != "" {
			activity.WeatherStatus = weather.WeatherDesc
		} else {
			activity.WeatherStatus = "Cuaca tidak diketahui (Belum Sync)"
		}
	}

	if err := database.DB.Save(&activity).Error; err != nil {
		log.Printf("Gagal update DB: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menyimpan perubahan ke database"})
		return
	}

	writeJSON(w, http.StatusOK, activity)
}

// DeleteActivity godoc
// @Summary      Hapus Kegiatan
// @Description  Menghapus data kegiatan berdasarkan ID
// @Tags         Activity
// @Param        id   query     string  true  "ID Kegiatan"
// @Success      200  {object}  map[string]string
// @Router       /api/activity/delete [delete]
func (h *ActivityHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "ID wajib diisi"})
		return
	}

	if err := database.DB.Delete(&models.Activity{}, id).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus kegiatan"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Kegiatan berhasil dihapus"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
