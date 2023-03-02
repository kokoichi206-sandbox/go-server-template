package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kokoichi206-sandbox/go-server-template/usecase"
	"kokoichi206-sandbox/go-server-template/util/logger"
)

type handler struct {
	logger  logger.Logger
	usecase usecase.Usecase

	Engine *gin.Engine
}

func New(logger logger.Logger, usecase usecase.Usecase) *handler {
	r := gin.Default()

	h := &handler{
		logger:  logger,
		usecase: usecase,
		Engine:  r,
	}
	h.setupRoutes()

	return h
}

func (h *handler) setupRoutes() {
	h.Engine.GET("/health", h.Health)
}

func (h *handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "ok",
	})
}
