package main

import (
	"context"
	"net"
	"os"

	"github.com/opentracing/opentracing-go"

	"kokoichi206-sandbox/go-server-template/config"
	"kokoichi206-sandbox/go-server-template/handler"
	"kokoichi206-sandbox/go-server-template/repository/database"
	"kokoichi206-sandbox/go-server-template/usecase"
	"kokoichi206-sandbox/go-server-template/util"
	"kokoichi206-sandbox/go-server-template/util/logger"
)

const (
	service = "server-template"
)

func main() {
	// config
	cfg := config.New()

	// logger
	logger := logger.NewBasicLogger(os.Stdout, "ubuntu", service)

	// tracer
	tracer, traceCloser, err := util.NewJaegerTracer(cfg.AgentHost, cfg.AgentPort, service)
	if err != nil {
		logger.Errorf(context.Background(), "cannot initialize jaeger tracer: ", err)
	} else {
		defer traceCloser.Close()
		opentracing.SetGlobalTracer(tracer)
	}

	// database
	database, err := database.New(
		cfg.DbDriver, cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword,
		cfg.DbName, cfg.DbSSLMode, logger,
	)
	if err != nil {
		logger.Errorf(context.Background(), "failed to db.New: ", err)
	}

	// usecase
	usecase := usecase.New(database, logger)

	// handler
	h := handler.New(logger, usecase)
	addr := net.JoinHostPort(cfg.ServerHost, cfg.ServerPort)

	// run
	if err := h.Engine.Run(addr); err != nil {
		logger.Critical(context.Background(), "failed to serve http")
	}
}
