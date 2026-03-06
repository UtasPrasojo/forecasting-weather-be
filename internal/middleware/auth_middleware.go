package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Gunakan secret key yang sama dengan saat Generate Token
var SecretKey = []byte("rahasia_telkom_sigma")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header diperlukan", http.StatusUnauthorized)
			return
		}

		var tokenString string

		// Cek apakah mengandung "Bearer "
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Jika tidak ada "Bearer", kita anggap seluruh header adalah token
			tokenString = authHeader
		}

		// Sekarang lanjutkan dengan parsing tokenString...
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ... (logika validasi tetap sama) ...
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak terduga: %v", token.Header["alg"])
			}
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token tidak valid atau kedaluwarsa", http.StatusUnauthorized)
			return
		}

		// (Opsional) Ambil data user dari claims dan masukkan ke context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Printf("Claims ditemukan: %v\n", claims)
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			next(w, r.WithContext(ctx))
		} else {
			next(w, r)
		}
	}
}
