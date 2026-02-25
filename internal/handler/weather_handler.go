package handler

import (
	"fmt"
	"net/http"
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/models"
	"project-telkom-sigma/internal/services"
	"time"
)

type WeatherHandler struct {
	WeatherService services.WeatherService
}

// NewWeatherHandler adalah Constructor untuk inisialisasi Handler
func NewWeatherHandler(s services.WeatherService) *WeatherHandler {
	return &WeatherHandler{WeatherService: s}
}

// HandleSync godoc
// @Summary      Sinkronisasi Data BMKG
// @Description  Mengambil data cuaca dari API BMKG berdasarkan kode wilayah (adm4) dan menyimpannya ke database
// @Tags         Weather
// @Param        adm4  query     string  true  "Kode Wilayah (contoh: 31.71.01.1001)"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/weather/sync [post]
func (h *WeatherHandler) HandleSync(w http.ResponseWriter, r *http.Request) {
	adm4 := r.URL.Query().Get("adm4")
	fmt.Println("adm4:", adm4)
	if adm4 == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "adm4 is required"})
		return
	}

	err := h.WeatherService.SyncWeather(adm4)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal Sync", "error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message":   "Sync Berhasil",
		"last_sync": time.Now().Format(time.RFC3339),
	})
}

// GetAllWeather godoc
// @Summary      Daftar Semua Data Cuaca
// @Description  Mengambil semua data cuaca yang telah di-sync, diurutkan berdasarkan waktu sinkronisasi terbaru
// @Tags         Weather
// @Produce      json
// @Success      200   {array}   models.Weather
// @Failure      500   {object}  map[string]string
// @Router       /api/weather [get]
func (h *WeatherHandler) GetAllWeather(w http.ResponseWriter, r *http.Request) {
	var weathers []models.Weather
	result := database.DB.Order("sync_time desc").Find(&weathers)
	if result.Error != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal ambil data"})
		return
	}
	writeJSON(w, http.StatusOK, weathers)
}

// GetWeatherByID godoc
// @Summary      Detail Data Cuaca
// @Tags         Weather
// @Param        id   query     string  true  "ID Weather"
// @Success      200  {object}  map[string]string
// @Router       /api/weather/detail [get]
func (h *WeatherHandler) GetWeatherByID(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Detail Data"})
}

// UpdateWeather godoc
// @Summary      Update Data Cuaca
// @Tags         Weather
// @Accept       json
// @Param        id      query     string          true  "ID Weather"
// @Param        weather body      models.Weather  true  "Data Weather"
// @Success      200     {object}  map[string]string
// @Router       /api/weather/update [put]
func (h *WeatherHandler) UpdateWeather(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Update Berhasil"})
}

// DeleteWeather godoc
// @Summary      Hapus Data Cuaca
// @Tags         Weather
// @Param        id   query     string  true  "ID Weather"
// @Success      200  {object}  map[string]string
// @Router       /api/weather/delete [delete]
func (h *WeatherHandler) DeleteWeather(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Delete Berhasil"})
}

// GetWeatherStats godoc
// @Summary      Data Statistik Dashboard
// @Description  Mengambil data agregasi untuk Pie Chart (Kategori) dan Column Chart (Suhu Harian)
// @Tags         Dashboard
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Contoh: {pie_chart: [], column_chart: []}"
// @Router       /api/weather/dashboard [get]
func (h *WeatherHandler) GetWeatherStats(w http.ResponseWriter, r *http.Request) {
	type PieData struct {
		Category string `json:"category"`
		Total    int64  `json:"total"`
	}
	var pieStats []PieData
	database.DB.Model(&models.Weather{}).Select("category, count(*) as total").Group("category").Scan(&pieStats)

	var columnStats []map[string]interface{}
	database.DB.Model(&models.Weather{}).
		Select("DATE(local_datetime) as day, AVG(t) as avg_temp").
		Group("DATE(local_datetime)").
		Order("day ASC").
		Scan(&columnStats)

	response := map[string]interface{}{
		"pie_chart":    pieStats,
		"column_chart": columnStats,
	}

	writeJSON(w, http.StatusOK, response)
}