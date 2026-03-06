package handler

import (
	"encoding/json"
	"net/http"
	"project-telkom-sigma/internal/database"
	"project-telkom-sigma/internal/middleware"
	"project-telkom-sigma/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary      User Register
// @Description  Mendaftarkan user baru ke dalam sistem
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "Data User (Username & Password)"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Router       /api/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal memproses password"})
		return
	}
	user.Password = string(hashedPassword)

	// Simpan ke database
	if err := database.DB.Create(&user).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Username sudah digunakan atau gagal simpan"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Registrasi berhasil"})
}

// Login godoc
// @Summary      User Login
// @Description  Autentikasi user dan mendapatkan JWT Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login  body      models.LoginRequest  true  "Data Login"
// @Success      200    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Router       /api/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginRequest // Menggunakan struct yang sudah didefinisikan

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
		return
	}

	var user models.User
	// Cari user berdasarkan username
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Username atau password salah"})
		return
	}

	// Bandingkan password input dengan hash di DB
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "Username atau password salah"})
		return
	}

	// Buat JWT Token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.SecretKey)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal generate token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": tokenString,
		"type":  "Bearer",
	})
}
