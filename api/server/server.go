package server

import (
	"net/http"
	"strconv"

	"github.com/estenity95/go-test-task/internal/config"
)

func New(handler http.Handler, cfg *config.AppConfig) *http.Server {
	return &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}
