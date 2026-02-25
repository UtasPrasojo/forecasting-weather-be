package handler

import (
	"encoding/json"
	"net/http"
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/models"
)

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
func CreateActivity(w http.ResponseWriter, r *http.Request) {
	var input models.Activity
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		return
	}

	var weather models.Weather
	err := database.DB.Where("area_code = ? AND local_datetime <= ?", input.AreaCode, input.ActivityDate).
		Order("local_datetime DESC").
		First(&weather).Error

	if err == nil {
		input.WeatherStatus = weather.WeatherDesc
	} else {
		input.WeatherStatus = "Cuaca tidak diketahui (Belum Sync)"
	}

	database.DB.Create(&input)
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
func GetAllActivities(w http.ResponseWriter, r *http.Request) {
	var activities []models.Activity

	query := r.URL.Query()
	search := query.Get("search")
	sortBy := query.Get("sort_by")
	order := query.Get("order")

	dbQuery := database.DB.Model(&models.Activity{})

	if search != "" {
		searchText := "%" + search + "%"
		dbQuery = dbQuery.Where("name ILIKE ? OR weather_status ILIKE ?", searchText, searchText)
	}

	if sortBy == "" {
		sortBy = "activity_date"
	}
	if order == "" || (order != "asc" && order != "desc") {
		order = "asc"
	}

	sortQuery := sortBy + " " + order

	err := dbQuery.Order(sortQuery).Find(&activities).Error
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil data kegiatan"})
		return
	}

	writeJSON(w, http.StatusOK, activities)
}

// DeleteActivity godoc
// @Summary      Hapus Kegiatan
// @Description  Menghapus data kegiatan berdasarkan ID
// @Tags         Activity
// @Param        id   query     string  true  "ID Kegiatan"
// @Success      200  {object}  map[string]string
// @Router       /api/activity/delete [delete]
func DeleteActivity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	database.DB.Delete(&models.Activity{}, id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Kegiatan dibatalkan"})
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
func UpdateActivity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "ID is required"})
		return
	}

	var activity models.Activity
	if err := database.DB.First(&activity, id).Error; err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "Kegiatan tidak ditemukan"})
		return
	}

	var input models.Activity
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		return
	}

	activity.Name = input.Name
	activity.AreaCode = input.AreaCode
	activity.ActivityDate = input.ActivityDate

	var weather models.Weather
	err := database.DB.Where("area_code = ? AND local_datetime <= ?", activity.AreaCode, activity.ActivityDate).
		Order("local_datetime DESC").
		First(&weather).Error

	if err == nil {
		activity.WeatherStatus = weather.WeatherDesc
	} else {
		activity.WeatherStatus = "Cuaca tidak diketahui (Belum Sync)"
	}

	database.DB.Save(&activity)
	writeJSON(w, http.StatusOK, activity)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
