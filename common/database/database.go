package database

import (
	"common/commonconfig"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

type Connection struct {
	*pgxpool.Pool
	ConnectionString string
}

type PgError struct {
	Code string
}

var (
	DB *Connection
)

func Connect() {
	if DB != nil {
		return
	}

	log.Printf("Connecting to database")
	log.Printf("Database host: %s", commonconfig.DbHost)
	dbConfig, connString := getConfig()
	conn, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
	}
	ping(conn)
	DB = &Connection{conn, connString}
}

func Close() {
	log.Printf("Disconnected from database")
	DB.Close()
}

func getConfig() (*pgxpool.Config, string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s", commonconfig.DbUsername, commonconfig.DbPassword, commonconfig.DbHost, commonconfig.DbPort, commonconfig.DbDatabase, commonconfig.DbSslMode, commonconfig.DbSchema)
	dbConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Printf("Error parsing database config: %v", err)
	}
	return dbConfig, connectionString
}

func ping(conn *pgxpool.Pool) {
	err := conn.Ping(context.Background())
	if err != nil {
		log.Printf("Error pinging database: %v", err)
	} else {
		log.Printf("Connected to database")
	}
}
