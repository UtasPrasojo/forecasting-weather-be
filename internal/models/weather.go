package models

import "time"

type Weather struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AreaCode      string    `gorm:"uniqueIndex:idx_weather_unique"`
	UtcDatetime   time.Time `gorm:"uniqueIndex:idx_weather_unique"`
	LocalDatetime time.Time `json:"local_datetime"`

	// Data Cuaca (sesuai parameter BMKG)
	T             float64 `json:"t"`               // Suhu
	Hu            int     `json:"hu"`              // Kelembapan
	WeatherDesc   string  `json:"weather_desc"`    // Kondisi (Indo)
	WeatherDescEn string  `json:"weather_desc_en"` // Kondisi (English)
	Ws            float64 `json:"ws"`              // Kecepatan Angin
	Wd            string  `json:"wd"`              // Arah Angin
	Tcc           int     `json:"tcc"`             // Tutupan Awan
	VsText        string  `json:"vs_text"`         // Jarak Pandang

	// Kriteria Tugas: Kategori & Sync Time
	Category     string    `gorm:"index" json:"category"` // Kita isi dari weather_desc (Cerah, Hujan, dll)
	AnalysisDate time.Time `json:"analysis_date"`
	SyncTime     time.Time `json:"sync_time"` // Kapan kita menembak API
}
type BMKGResponse struct {
	Data []struct {
		Lokasi struct {
			Adm4 string `json:"adm4"`
		} `json:"lokasi"`
		Cuaca [][]struct {
			UtcDatetime   string  `json:"utc_datetime"`
			LocalDatetime string  `json:"local_datetime"`
			T             float64 `json:"t"`
			Hu            int     `json:"hu"`
			WeatherDesc   string  `json:"weather_desc"`
		} `json:"cuaca"`
	} `json:"data"`
}
