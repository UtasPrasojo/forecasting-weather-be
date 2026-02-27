package routes

import (
	"net/http"
	_ "project-telkom-sigma/docs"
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/handler"
	"project-telkom-sigma/internal/repositories"
	"project-telkom-sigma/internal/services"

	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

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

	activityHandler := handler.NewActivityHandler()

	mux.HandleFunc("POST /api/activity", activityHandler.CreateActivity)
	mux.HandleFunc("GET /api/activity", activityHandler.GetAllActivities)
	mux.HandleFunc("PUT /api/activity/update", activityHandler.UpdateActivity)
	mux.HandleFunc("DELETE /api/activity/delete", activityHandler.DeleteActivity)


	wilayahRepo := repositories.NewWilayahRepository(database.DB)
	wilayahService := services.NewWilayahService(wilayahRepo)
	wilayahHandler := handler.NewWilayahHandler(wilayahService)

	// Route Wilayah
	mux.HandleFunc("GET /api/wilayah", wilayahHandler.GetWilayah)
	return c.Handler(mux)
}
