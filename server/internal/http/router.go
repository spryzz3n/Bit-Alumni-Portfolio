package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Store is the data access the handlers need; kept small so tests can fake it.
type Store interface {
	AppName(ctx context.Context) (string, error)
}

type Server struct {
	store  Store
	logger *slog.Logger
}

func NewRouter(store Store, logger *slog.Logger, corsOrigin string) http.Handler {
	s := &Server{store: store, logger: logger}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(s.logRequests)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{corsOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		MaxAge:         300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", s.health)
		r.Get("/db-check", s.dbCheck)
	})
	return r
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) dbCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	name, err := s.store.AppName(ctx)
	if err != nil {
		s.logger.Error("db check failed", "err", err, "request_id", middleware.GetReqID(r.Context()))
		writeError(w, http.StatusServiceUnavailable, "db_unavailable", "The database is not available right now.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "app_name": name})
}

func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		next.ServeHTTP(ww, r)
		s.logger.Info("request",
			"method", r.Method, "path", r.URL.Path, "status", ww.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", middleware.GetReqID(r.Context()))
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
