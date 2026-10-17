package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPAddr      string
	CORSOrigin    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	MigrationsDir string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		CORSOrigin:    getenv("CORS_ORIGIN", "http://localhost:5173"),
		DBHost:        getenv("DB_HOST", "127.0.0.1"),
		DBPort:        getenv("DB_PORT", "3306"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		MigrationsDir: getenv("MIGRATIONS_DIR", "../migrations"),
	}
	if cfg.DBUser == "" || cfg.DBName == "" {
		return cfg, fmt.Errorf("DB_USER and DB_NAME must be set")
	}
	return cfg, nil
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
