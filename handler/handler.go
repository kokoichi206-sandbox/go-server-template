package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kokoichi206-sandbox/go-server-template/model/apperr"
	"github.com/kokoichi206-sandbox/go-server-template/usecase"
	"github.com/kokoichi206-sandbox/go-server-template/util/logger"
)

type handler struct {
	logger  logger.Logger
	usecase usecase.Usecase

	Engine *gin.Engine
}

//nolint:revive
func New(logger logger.Logger, usecase usecase.Usecase) *handler {
	r := gin.New()
	r.Use(requestLog())
	r.Use(gin.Recovery())

	h := &handler{
		logger:  logger,
		usecase: usecase,
		Engine:  r,
	}
	h.setupRoutes()

	return h
}

type ginLog struct {
	TimeStamp  time.Time     `json:"time_stamp"`
	StatusCode int           `json:"status"`
	Latency    time.Duration `json:"latency"`
	ClientIP   string        `json:"client_ip"`
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	RequestID  string        `json:"request_id"`
}

func requestLog() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: formatter,
		SkipPaths: []string{},
	})
}

func formatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	gl := ginLog{
		TimeStamp:  param.TimeStamp,
		StatusCode: param.StatusCode,
		Latency:    param.Latency.Truncate(time.Millisecond),
		ClientIP:   param.ClientIP,
		Method:     param.Method,
		Path:       param.Path,
	}

	//nolint:errcheck
	b, _ := json.Marshal(gl)

	return fmt.Sprintln(string(b))
}

func (h *handler) setupRoutes() {
	base := h.Engine.Group("/api/v1")
	base.Use(h.requestIDMW())

	base.Handle(http.MethodGet, "/health", handlerWrapper(h.Health, h.logger))
}

func handlerWrapper(fun func(c *gin.Context) error, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fun(c); err != nil {
			handleError(c, logger, err)
		}
	}
}

// handleError is a helper function to handle error.
// This function writes status code and error message to response body.
func handleError(c *gin.Context, logger logger.Logger, err error) {
	var e apperr.AppError
	if ok := errors.As(err, &e); !ok {
		e = apperr.AppError{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal server error",
			Log:        err.Error(),
		}
	}

	if e.Log != "" {
		logger.Error(context.WithoutCancel(c.Request.Context()), e.Log)
	}

	c.JSON(e.StatusCode, gin.H{
		"error": e.Message,
	})
}

func (h *handler) Health(c *gin.Context) error {
	c.JSON(http.StatusOK, gin.H{
		"health": "ok",
	})

	return nil
}
