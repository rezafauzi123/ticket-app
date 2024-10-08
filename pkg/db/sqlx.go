package db

import (
	"fmt"
	"ticket-app/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB(cfg *config.DBConfig) (*sqlx.DB, error) {
	dbSource := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

	db, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}
