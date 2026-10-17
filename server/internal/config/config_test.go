package config

import "testing"

func TestLoadRequiresDBVars(t *testing.T) {
	t.Setenv("DB_USER", "")
	t.Setenv("DB_NAME", "")
	if _, err := Load(); err == nil {
		t.Fatal("expected error when DB_USER/DB_NAME are unset")
	}
}

func TestLoadDefaultsAndDSN(t *testing.T) {
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_PASSWORD", "p")
	t.Setenv("DB_NAME", "n")
	t.Setenv("DB_HOST", "")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := "u:p@tcp(127.0.0.1:3306)/n?parseTime=true&charset=utf8mb4"
	if got := cfg.DSN(); got != want {
		t.Fatalf("DSN = %q, want %q", got, want)
	}
}
