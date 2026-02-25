package routes

import (
	"net/http"
	_ "project-telkom-sigma/docs"
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/handler" 
	"project-telkom-sigma/internal/repositories"
	"project-telkom-sigma/internal/services"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	weatherRepo := repositories.NewWeatherRepository(database.DB)

	weatherService := services.NewWeatherService(weatherRepo)

	weatherHandler := handler.NewWeatherHandler(weatherService)

	mux.HandleFunc("POST /api/weather/sync", weatherHandler.HandleSync)
	mux.HandleFunc("GET /api/weather", weatherHandler.GetAllWeather)
	mux.HandleFunc("GET /api/weather/detail", weatherHandler.GetWeatherByID)
	mux.HandleFunc("PUT /api/weather/update", weatherHandler.UpdateWeather)
	mux.HandleFunc("DELETE /api/weather/delete", weatherHandler.DeleteWeather)
	mux.HandleFunc("GET /api/weather/dashboard", weatherHandler.GetWeatherStats)

	return mux
}
