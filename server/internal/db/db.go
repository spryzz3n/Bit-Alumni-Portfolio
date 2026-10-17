package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Open(dsn string) (*sql.DB, error) {
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(10)
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(30 * time.Minute)
	return d, nil
}

// AppName reads the seed row created by migration 0001.
func AppName(ctx context.Context, d *sql.DB) (string, error) {
	var v string
	err := d.QueryRowContext(ctx,
		"SELECT info_value FROM app_info WHERE info_key = ?", "app_name").Scan(&v)
	return v, err
}
