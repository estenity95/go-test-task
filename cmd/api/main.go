package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/estenity95/go-test-task/api/docs"
	"github.com/estenity95/go-test-task/api/resource/subscription"
	"github.com/estenity95/go-test-task/api/router"
	"github.com/estenity95/go-test-task/api/server"
	"github.com/estenity95/go-test-task/internal/config"
	"github.com/estenity95/go-test-task/internal/logger"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

// @title           Subscriptions API
// @version         1.0
// @description     REST API для управления подписками.
// @basePath        /api/v1
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New(cfg.DB.Debug)
	validator := validator.New()

	dsn := fmt.Sprintf(fmtDBString, cfg.DB.Host, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Port)

	logLevel := gormlogger.Warn
	if cfg.DB.Debug {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		logger.Fatal().Err(err).Msg("DB connection start failure")
		return
	}

	repo := subscription.NewRepo(db)
	r := router.NewRouter(logger, validator, repo)
	server := server.New(r, cfg)

	go func() {
		logger.Info().Msg("http server starting")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			logger.Error().Err(err).Msg("http server failed")
			os.Exit(1)
		}
	}()

	// ожидаем SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
	logger.Info().Msg("shutdown complete")
}
