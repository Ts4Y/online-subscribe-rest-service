package main

import (
	"context"
	"fmt"
	"net/http"
	"online-subscribe-rest-service/internal/api/handler"
	"online-subscribe-rest-service/internal/api/router"
	"online-subscribe-rest-service/internal/repository"
	"online-subscribe-rest-service/internal/service"
	"online-subscribe-rest-service/pkg/config"
	"online-subscribe-rest-service/pkg/logger"
	"online-subscribe-rest-service/pkg/postgres"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New(".env")
	if err != nil {
		fmt.Printf("failed to load config: %s\n", err.Error())
		return
	}

	log, err := logger.New(cfg.Logger.Mode)
	if err != nil {
		fmt.Printf("failed to create logger: %s\n", err.Error())
		return
	}

	pgConn, err := postgres.ConnectToPostgres(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.ErrorF("failed to connect to postgres: %w", err)
		return
	}

	if err := postgres.UpMigrations(cfg.Postgres.DSN); err != nil {
		log.ErrorF("failed to up migrations: %w", err)
		return
	}

	repo := repository.NewSubscriptionRepo(pgConn)
	service := service.NewService(repo)
	handler := handler.NewHandler(log, service)
	router := router.NewRouter(handler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		Handler:      router,
	}

	log.InfoF("server succesfully started on port %d", cfg.HTTP.Port)

	if err := server.ListenAndServe(); err != nil {
		log.ErrorF("failed to run http server: %w", err)
	}
}
