package database

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"project-telkom-sigma/internal/configs"
	"project-telkom-sigma/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

func InitDB(cfg *configs.Setting) {
	dsn := cfg.Database.ConnStr

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	err = db.AutoMigrate(&models.Wilayah{})
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}

	DB = db
	fmt.Println("Database terkoneksi!")
}

func SeedWilayah(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Peringatan: Gagal membuka file CSV: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		log.Println("Gagal membaca header CSV")
		return
	}

	fmt.Println("Memulai proses input data wilayah...")

	// ... di dalam loop for di SeedWilayah ...

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		// Mapping sesuai isi CSV kamu:
		// record[0] adalah 'code' (misal: 11.01)
		// record[1] adalah 'Loc' (misal: KAB. ACEH SELATAN)
		wilayah := models.Wilayah{
			Code: record[0],
			Loc:  record[1],
		}

		// Simpan atau Update jika data sudah ada
		DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}}, // Sesuaikan dengan primary key
			DoUpdates: clause.AssignmentColumns([]string{"loc", "updated_at"}),
		}).Create(&wilayah)
	}

	fmt.Println("Proses input CSV selesai!")
}
