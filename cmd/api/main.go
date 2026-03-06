package main

import (
	"log"
	"net/http"
	_ "project-telkom-sigma/docs"
	"project-telkom-sigma/internal/configs"
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/routes"
)

// @title           API Sistem Informasi Telkom Sigma
// @version         1.0
// @description     Dokumentasi API untuk Proyek Internal
// @host            localhost:8080
// @BasePath        /api

// Tambahkan blok di bawah ini:
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @scheme bearer
// @bearerFormat JWT
func main() {
	setting, err := configs.NewSetting()
	if err != nil {
		log.Fatal("Gagal load config:", err)
	}

	database.InitDB(setting)

	router := routes.SetupRoutes()

	log.Println("Server mulai di port :8080")
	log.Println("Dokumentasi Swagger: http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
