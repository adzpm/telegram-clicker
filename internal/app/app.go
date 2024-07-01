package app

import (
	zap "go.uber.org/zap"

	config "github.com/adzpm/telegram-clicker/internal/config"
	math "github.com/adzpm/telegram-clicker/internal/math"
	rest "github.com/adzpm/telegram-clicker/internal/rest"
	storage "github.com/adzpm/telegram-clicker/internal/storage"
)

type (
	App struct {
		rest    *rest.REST
		config  *config.Config
		storage *storage.Storage
		math    *math.Math
		logger  *zap.Logger
	}
)

func New(
	r *rest.REST,
	c *config.Config,
	s *storage.Storage,
	m *math.Math,
	l *zap.Logger,
) *App {
	return &App{
		rest:    r,
		config:  c,
		storage: s,
		math:    m,
		logger:  l,
	}
}
