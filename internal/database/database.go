package database

import (
	"fmt"
	"log"
	"project-telkom-sigma/internal/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *configs.Setting) {
	dsn := cfg.Database.ConnStr

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	DB = db
	fmt.Println("Database terkoneksi!")
}
