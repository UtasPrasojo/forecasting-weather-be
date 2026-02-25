package main

import (
	"log"
	"net/http"
	"project-telkom-sigma/internal/configs"
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/routes"
	
)

// @title           API Analitik Publik
// @version         1.0
// @description     Backend service untuk konsumsi API dan Dashboard Analitik.
// @host            localhost:8080
// @BasePath        /
func main() {

	setting, err := configs.NewSetting()
	if err != nil {
		log.Fatal("Gagal load config:", err)
	}

	// ✅ INIT DATABASE DULU
	database.InitDB(setting)
	database.SeedWilayah("./internal/migration/wilayah.csv")

	router := routes.NewRouter()

	log.Println("Server mulai di port :8080")
	log.Println("Dokumentasi Swagger dapat diakses di: http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
