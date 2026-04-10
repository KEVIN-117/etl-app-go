package main

import (
	"etl_pipeline/pkg/config"
	"etl_pipeline/pkg/db"
	"etl_pipeline/pkg/logger"
	"net/http"
)

func main() {
	logger.Info("Hello World")
	cfg := config.LoadConfig()
	db := db.Connect(cfg)
	// defer db.Close()

	logger.Info("Server running on :8080", logger.F("db", db != nil))

	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API running 🚀"))
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("Server failed to start", logger.F("error", err))
	}
}
