package main

import (
	_ "EMTT/docs"
	"EMTT/internal/config"
	"EMTT/internal/handlers"
	"EMTT/internal/repository"
	handler "EMTT/internal/server"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// @title Subscription Service API
// @version 1.0
// @description API for managing subscriptions
// @host localhost:8088
// @BasePath /
func main() {
	ctx := context.Background()
	log := slog.Default()

	cfg := config.Load()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB_USER, cfg.DB_PASS, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)

	log.Info("Connecting to database",
		"host", cfg.DB_HOST,
		"port", cfg.DB_PORT,
		"database", cfg.DB_NAME)

	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Error("Database connection error", "error", err)
		return
	}
	defer db.Close()

	if err = db.Ping(ctx); err != nil {
		log.Error("Database ping failed", "error", err)
		return
	}

	repo := repository.NewSubScriptRepo(db, ctx)
	h := &handlers.Handler{SubRepo: repo}

	server := handler.NewServer(":"+cfg.SER_PORT, log)

	log.Info("Server starting on port " + cfg.SER_PORT)
	server.Start(h)
}
