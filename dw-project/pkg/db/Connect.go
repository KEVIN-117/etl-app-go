package db

import (
	"database/sql"
	"etl_pipeline/pkg/config"
	"etl_pipeline/pkg/logger"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to connect to database", logger.F("error", err))
	}

	if err = db.Ping(); err != nil {
		logger.Error("Failed to ping database", logger.F("error", err))
	}

	logger.Info("Successfully connected to database")
	return db
}
