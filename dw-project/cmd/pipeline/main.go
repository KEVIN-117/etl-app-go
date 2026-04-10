package main

import (
	"etl_pipeline/pkg/config"
	"etl_pipeline/pkg/db"
	"etl_pipeline/pkg/logger"
)

func main() {
	logger.Info("Hello World")
	cfg := config.LoadConfig()
	database := db.Connect(cfg)

	logger.Info("ETL listo con DB", logger.F("db", database != nil))
}
