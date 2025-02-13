package database

import (
	"fmt"
	"log"
	"sync"

	"Backend/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(wg *sync.WaitGroup) {
	defer wg.Done()

	cfg := config.GetConfig();
	dbURL := cfg.DBUrl;
	
	m, err := migrate.New(
		"file://database/migrations",
		dbURL,
	)		

	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Println("[LOG]: Migrations Applied Successfully !");
}
