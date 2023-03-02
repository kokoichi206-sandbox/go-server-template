package usecase

import (
	"kokoichi206-sandbox/go-server-template/repository"
	"kokoichi206-sandbox/go-server-template/util/logger"
)

type Usecase interface {
}

type usecase struct {
	database repository.Database

	logger logger.Logger
}

func New(database repository.Database, logger logger.Logger) Usecase {
	usecase := &usecase{
		database: database,
		logger:   logger,
	}

	return usecase
}
