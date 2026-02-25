package routes

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "project-telkom-sigma/docs" // ⚠️ WAJIB untuk register swagger
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Swagger endpoint
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// contoh endpoint lain
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}