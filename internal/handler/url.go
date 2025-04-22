package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLService interface {
	Shorten(original string) (string, error)
	Resolve(short string) (string, error)
}

type URLHandler struct {
	service URLService
}

func NewURLHandler(service URLService) *URLHandler {
	return &URLHandler{service: service}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"url"`
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWirhError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	short, err := h.service.Shorten(req.URL)
	if err != nil {
		respondWirhError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp := ShortenResponse{ShortURL: "http://localhost:8080/" + short}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *URLHandler) ResolveURL(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "shortUrl")
	originalUrl, err := h.service.Resolve(short)
	if err != nil {
		respondWirhError(w, http.StatusNotFound, "Not found")
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)
}

func respondWirhError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
