package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/kokoichi206-sandbox/go-server-template/config"
	"github.com/kokoichi206-sandbox/go-server-template/handler"
	"github.com/kokoichi206-sandbox/go-server-template/repository/database"
	"github.com/kokoichi206-sandbox/go-server-template/usecase"
	"github.com/kokoichi206-sandbox/go-server-template/util"
	"github.com/kokoichi206-sandbox/go-server-template/util/logger"
)

const (
	service           = "server-template"
	gracePeriod       = 5 * time.Second
	readHeaderTimeout = 5 * time.Second
)

func main() {
	// config
	cfg := config.New()

	// logger
	logger := logger.NewBasicLogger(os.Stdout, "ubuntu", service)

	// tracer
	tracer, traceCloser, err := util.NewJaegerTracer(cfg.AgentHost, cfg.AgentPort, service)
	if err != nil {
		logger.Errorf(context.Background(), "cannot initialize jaeger tracer: %s", err)
	} else {
		defer traceCloser.Close()
		opentracing.SetGlobalTracer(tracer)
	}

	// database
	database, err := database.New(
		cfg.DBDriver, cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword,
		cfg.DBName, cfg.DBSSLMode, logger,
	)
	if err != nil {
		logger.Errorf(context.Background(), "failed to db.New: %s", err)
	}

	// usecase
	usecase := usecase.New(database, logger)

	// handler
	h := handler.New(logger, usecase)
	addr := net.JoinHostPort(cfg.ServerHost, cfg.ServerPort)

	// run
	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Engine,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Errorf(context.Background(), "failed to listen and serve: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(context.Background(), "Server is shutting down...")

	// 5 seconds grace period.
	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf(context.Background(), "failed to server shutdown: %s", err)
	}

	logger.Info(context.Background(), "Server is gracefully stopped")
}
