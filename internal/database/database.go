package database

import (
    "fmt"
    "log"
    "gorm.io/driver/postgres" 
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    dsn := "host=localhost user=user password=pass dbname=apotek port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Gagal koneksi database:", err)
    }
    
    DB = db
    fmt.Println("Database terkoneksi!")
}