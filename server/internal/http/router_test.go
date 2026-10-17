package http

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeStore struct{ err error }

func (f fakeStore) AppName(context.Context) (string, error) { return "Alumni Portal", f.err }

func newTestRouter(s Store) http.Handler {
	return NewRouter(s, slog.New(slog.NewTextHandler(io.Discard, nil)), "http://localhost:5173")
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestRouter(fakeStore{err: errors.New("down")}).ServeHTTP(rec, httptest.NewRequest("GET", "/api/v1/health", nil))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"ok"`) {
		t.Fatalf("got %d %s", rec.Code, rec.Body)
	}
}

func TestDBCheckOK(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestRouter(fakeStore{}).ServeHTTP(rec, httptest.NewRequest("GET", "/api/v1/db-check", nil))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Alumni Portal") {
		t.Fatalf("got %d %s", rec.Code, rec.Body)
	}
}

func TestDBCheckHidesDBError(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestRouter(fakeStore{err: errors.New("secret sql detail")}).ServeHTTP(rec, httptest.NewRequest("GET", "/api/v1/db-check", nil))
	if rec.Code != 503 || strings.Contains(rec.Body.String(), "secret") {
		t.Fatalf("got %d %s", rec.Code, rec.Body)
	}
}
