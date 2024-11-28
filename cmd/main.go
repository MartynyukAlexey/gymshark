package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/MartynyukAlexey/gymshark/internal/api"
	"github.com/MartynyukAlexey/gymshark/internal/config"
	"github.com/MartynyukAlexey/gymshark/internal/service"
	"github.com/MartynyukAlexey/gymshark/internal/smtp"
	"github.com/MartynyukAlexey/gymshark/internal/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	config := config.GetConfig()

	postgres, err := openPostgres(config.Postgres)
	if err != nil {
		logger.Error("db connection pool startup error", "err", err.Error())
		os.Exit(-1)
	}

	defer func() {
		if err := postgres.Close(); err != nil {
			logger.Error("db connection pool shutdown error", "err", err.Error())
		}
		logger.Info("db connection pool shutdown")
	}()

	minioClient, err := openMinio(config.Minio)
	if err != nil {
		logger.Error("minio startup error", "err", err.Error())
		os.Exit(-1)
	}

	store := storage.NewStorage(postgres, minioClient)

	mailer := smtp.NewSMTPMailer(config.Mailer, logger)

	svc := service.NewService(&service.ServiceOpts{
		Storage:    store,
		Mailer:     mailer,
		Logger:     logger,
		AuthConfig: config.Auth,
	})

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Server.Port),
		WriteTimeout: config.Server.WriteTimeout,
		ReadTimeout:  config.Server.ReadTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
		Handler:      api.Routes(svc, logger),
	}

	go func() {
		logger.Info("starting serving new connections...")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", "err", err)
			os.Exit(-1)
		}
		logger.Info("stopped serving new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "err", err)
	}
	if err := postgres.Close(); err != nil {
		logger.Error("db connection pool shutdown error", "err", err)
	}

	logger.Info("server shutdown")
}

func openPostgres(config *config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(config.MaxIdleTime)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func openMinio(config *config.MinioConfig) (*minio.Client, error) {
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(config.User, config.Password, ""),
	}

	minioClient, err := minio.New(config.Endpoint, opts)
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
