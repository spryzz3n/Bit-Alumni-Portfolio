// Command migrate applies versioned SQL migrations: `go run ./cmd/migrate up|down`.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/spryzz3n/bitalumni/server/internal/config"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		fmt.Fprintln(os.Stderr, "usage: migrate up|down")
		os.Exit(2)
	}
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	dir, err := filepath.Abs(cfg.MigrationsDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	m, err := migrate.New("file://"+filepath.ToSlash(dir), "mysql://"+cfg.DSN())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer m.Close()

	if os.Args[1] == "up" {
		err = m.Up()
	} else {
		err = m.Steps(-1)
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("migrate", os.Args[1], "done")
}
