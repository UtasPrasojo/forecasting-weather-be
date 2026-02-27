package models

import "time"

type Weather struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	AreaCode string `gorm:"uniqueIndex:idx_weather_unique" json:"area_code"`

	Wilayah Wilayah `gorm:"foreignKey:AreaCode;references:Code" json:"wilayah"`

	UtcDatetime   time.Time `gorm:"uniqueIndex:idx_weather_unique" json:"utc_datetime"`
	LocalDatetime time.Time `json:"local_datetime"`

	T             float64 `json:"t"`
	Hu            int     `json:"hu"`
	WeatherDesc   string  `json:"weather_desc"`
	WeatherDescEn string  `json:"weather_desc_en"`
	Ws            float64 `json:"ws"`
	Wd            string  `json:"wd"`
	Tcc           int     `json:"tcc"`
	VsText        string  `json:"vs_text"`

	Category     string    `gorm:"index" json:"category"`
	AnalysisDate time.Time `json:"analysis_date"`
	SyncTime     time.Time `json:"sync_time"`
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
