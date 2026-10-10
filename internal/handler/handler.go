package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"portfolio-marketplace/backend/internal/config"
	"portfolio-marketplace/backend/internal/middleware"
	"portfolio-marketplace/backend/internal/model"
	"portfolio-marketplace/backend/internal/repository"
)

type Handler struct {
	Repo   *repository.Postgres
	Config config.Config
}

func New(repo *repository.Postgres, cfg config.Config) *Handler {
	return &Handler{Repo: repo, Config: cfg}
}
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(h.cors)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, map[string]string{"status": "ok"}) })
	r.Post("/api/v1/auth/register", h.register)
	r.Post("/api/v1/auth/login", h.login)
	return r
}
func (h *Handler) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == h.Config.FrontendOrigin || origin == "http://localhost:5000" || origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var in model.RegisterRequest
	if !decode(w, r, &in) {
		return
	}
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.DisplayName = strings.TrimSpace(in.DisplayName)
	if !strings.Contains(in.Email, "@") || len(in.Password) < 8 || in.DisplayName == "" {
		writeJSON(w, 400, map[string]string{"error": "valid email, display name, and password of at least 8 characters are required"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		serverError(w, err)
		return
	}
	user, err := h.Repo.CreateUser(r.Context(), in.Email, string(hash), in.DisplayName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJSON(w, 409, map[string]string{"error": "email already registered"})
			return
		}
		serverError(w, err)
		return
	}
	token, err := middleware.CreateToken(user.ID, h.Config.JWTSecret)
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, 201, model.AuthResponse{Token: token, User: user})
}
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var in model.LoginRequest
	if !decode(w, r, &in) {
		return
	}
	user, hash, err := h.Repo.UserByEmail(r.Context(), strings.TrimSpace(strings.ToLower(in.Email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Password)) != nil {
		writeJSON(w, 401, map[string]string{"error": "invalid email or password"})
		return
	}
	token, err := middleware.CreateToken(user.ID, h.Config.JWTSecret)
	if err != nil {
		serverError(w, err)
		return
	}
	writeJSON(w, 200, model.AuthResponse{Token: token, User: user})
}
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	defer r.Body.Close()
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON request"})
		return false
	}
	return true
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}
func serverError(w http.ResponseWriter, err error) {
	log.Printf("request failed: %v", err)
	writeJSON(w, 500, map[string]string{"error": "internal server error"})
}
