package handler

import (
	"net/http"
	"project-telkom-sigma/internal/services"
)

type WilayahHandler struct {
	Service services.WilayahService
}

func NewWilayahHandler(s services.WilayahService) *WilayahHandler {
	return &WilayahHandler{Service: s}
}

// GetWilayah godoc
// @Summary      Daftar Wilayah (ADM4)
// @Description  Mengambil data wilayah untuk dropdown search
// @Tags         Wilayah
// @Param        q    query     string  false  "Cari kode atau nama wilayah"
// @Success      200  {array}   models.Wilayah
// @Router       /api/wilayah [get]
func (h *WilayahHandler) GetWilayah(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	data, err := h.Service.SearchWilayah(query)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil data wilayah"})
		return
	}
	writeJSON(w, http.StatusOK, data)
}