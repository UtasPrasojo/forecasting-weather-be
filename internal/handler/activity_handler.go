package handler

import (
	"encoding/json"
	"fmt"
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

// Helper untuk mengambil UserID dari Context
func getUserID(r *http.Request) uint {
	val := r.Context().Value("user_id")
	if val == nil {
		return 0
	}
	// JWT claims biasanya berupa float64 saat di-decode ke map
	if id, ok := val.(float64); ok {
		return uint(id)
	}
	return val.(uint)
}

// CreateActivity godoc
// @Summary      Buat Rencana Kegiatan
// @Description  Membuat rencana kegiatan baru milik user yang sedang login
// @Tags         Activity
// @Accept       json
// @Produce      json
// @Param        activity  body      models.Activity  true  "Data JSON Kegiatan"
// @Success      201       {object}  models.Activity
// @Failure      400       {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/activity [post]
func (h *ActivityHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	var input models.Activity
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
		return
	}

	// Set UserID dari token
	input.UserID = userID

	var lastWeather models.Weather
	if err := database.DB.Order("sync_time DESC").First(&lastWeather).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Belum ada data wilayah yang di-sync.",
		})
		return
	}

	if input.AreaCode == "" {
		input.AreaCode = lastWeather.AreaCode
	}

	// ... (Logika findWeather tetap sama) ...
	var weather models.Weather
	findWeather := func(code string, date time.Time) error {
		return database.DB.Where("area_code = ?", code).
			Order(fmt.Sprintf("ABS(EXTRACT(EPOCH FROM (local_datetime - '%s')))",
				date.Format("2006-01-02 15:04:05"))).
			First(&weather).Error
	}

	if err := findWeather(input.AreaCode, input.ActivityDate); err != nil {
		h.WeatherService.SyncWeather(input.AreaCode)
		findWeather(input.AreaCode, input.ActivityDate)
	}
	input.WeatherStatus = weather.WeatherDesc

	if err := database.DB.Create(&input).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal simpan kegiatan"})
		return
	}

	writeJSON(w, http.StatusCreated, input)
}

// GetAllActivities godoc
// @Summary      Daftar Semua Kegiatan User
// @Description  Mengambil daftar kegiatan milik user yang sedang login
// @Tags         Activity
// @Produce      json
// @Security     BearerAuth
// @Router       /api/activity [get]
func (h *ActivityHandler) GetAllActivities(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var activities []models.Activity

	// Filter berdasarkan UserID
	err := database.DB.Model(&models.Activity{}).
		Where("user_id = ?", userID).
		Preload("Wilayah").
		Find(&activities).Error

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil data"})
		return
	}
	writeJSON(w, http.StatusOK, activities)
}

// UpdateActivity godoc
// @Summary      Perbarui Kegiatan
// @Description  Mengubah data kegiatan milik user yang sedang login
// @Tags         Activity
// @Security     BearerAuth
// @Router       /api/activity/update [put]
func (h *ActivityHandler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	id := r.URL.Query().Get("id")

	var activity models.Activity
	// Filter dengan UserID agar user tidak bisa update milik orang lain
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&activity).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "Kegiatan tidak ditemukan atau bukan milik Anda"})
		return
	}

	// ... (Sisa logika update tetap sama) ...
	writeJSON(w, http.StatusOK, activity)
}

// DeleteActivity godoc
// @Summary      Hapus Kegiatan
// @Description  Menghapus kegiatan milik user yang sedang login
// @Tags         Activity
// @Security     BearerAuth
// @Router       /api/activity/delete [delete]
func (h *ActivityHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	id := r.URL.Query().Get("id")

	// Filter dengan UserID agar aman
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Activity{})
	if result.Error != nil || result.RowsAffected == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "Kegiatan tidak ditemukan atau akses ditolak"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Kegiatan berhasil dihapus"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
