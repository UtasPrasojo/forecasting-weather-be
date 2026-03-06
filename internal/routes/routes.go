package routes

import (
	"net/http"
	_ "project-telkom-sigma/docs" // Import wajib untuk register SwaggerInfo
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/handler"
	"project-telkom-sigma/internal/middleware"
	"project-telkom-sigma/internal/repositories"
	"project-telkom-sigma/internal/services"

	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// 1. Swagger UI Handler
	// Kita gunakan WrapHandler untuk melayani UI dan aset statisnya (JS/CSS)
	// Kita tambahkan URL eksplisit agar ia tahu harus mencari file di mana
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// 2. Redirect/Serve file swagger.json sebagai doc.json
	// Library Swagger mencari 'doc.json', tapi kita punya 'swagger.json'.
	// Kita buat mapping agar permintaan ke doc.json dilayani oleh file swagger.json Anda.
	mux.Handle("/swagger/doc.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	}))

	// 3. Konfigurasi CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// 4. API Mux (Terpisah agar tidak konflik dengan Swagger)
	apiMux := http.NewServeMux()

	// Inisialisasi Layer
	weatherRepo := repositories.NewWeatherRepository(database.DB)
	weatherService := services.NewWeatherService(weatherRepo)
	weatherHandler := handler.NewWeatherHandler(weatherService)
	activityHandler := handler.NewActivityHandler(weatherService)
	wilayahRepo := repositories.NewWilayahRepository(database.DB)
	wilayahService := services.NewWilayahService(wilayahRepo)
	wilayahHandler := handler.NewWilayahHandler(wilayahService)

	// Auth Routes
	apiMux.HandleFunc("POST /api/register", handler.Register)
	apiMux.HandleFunc("POST /api/login", handler.Login)

	// Weather Routes
	apiMux.HandleFunc("POST /api/weather/sync", middleware.AuthMiddleware(weatherHandler.HandleSync))
	apiMux.HandleFunc("GET /api/weather", weatherHandler.GetAllWeather)
	apiMux.HandleFunc("GET /api/weather/detail", weatherHandler.GetWeatherByID)
	apiMux.HandleFunc("PUT /api/weather/update", middleware.AuthMiddleware(weatherHandler.UpdateWeather))
	apiMux.HandleFunc("DELETE /api/weather/delete", middleware.AuthMiddleware(weatherHandler.DeleteWeather))
	apiMux.HandleFunc("GET /api/weather/dashboard", weatherHandler.GetWeatherStats)

	// Activity Routes
	apiMux.HandleFunc("POST /api/activity", middleware.AuthMiddleware(activityHandler.CreateActivity))
	apiMux.HandleFunc("GET /api/activity", middleware.AuthMiddleware(activityHandler.GetAllActivities))
	apiMux.HandleFunc("PUT /api/activity/update", middleware.AuthMiddleware(activityHandler.UpdateActivity))
	apiMux.HandleFunc("DELETE /api/activity/delete", middleware.AuthMiddleware(activityHandler.DeleteActivity))

	// Wilayah Routes
	apiMux.HandleFunc("GET /api/wilayah", wilayahHandler.GetWilayah)

	// Masukkan apiMux ke dalam filter CORS
	mux.Handle("/api/", c.Handler(apiMux))

	return mux
}
