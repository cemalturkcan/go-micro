package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"log"
)

func MigrateDb() {
	m, err := migrate.New(
		"file://migrations",
		DB.ConnectionString)

	if err != nil {
		log.Printf("Error creating migration: %v", err)
	}
	err = m.Up()
	if err != nil {
		log.Printf("Error running migration: %v", err)
	}
	log.Println("Migration completed")
}
