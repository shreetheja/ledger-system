package pg

import (
	"database/sql"
	"fmt"
	"log"

	"inoscipta/config"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var DB *sql.DB

func InitPostgres() {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDB,
		config.PostgresHost,
		config.PostgresPort,
	)

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open Postgres connection: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL")
}
