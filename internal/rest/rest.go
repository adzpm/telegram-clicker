package rest

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"
	zap "go.uber.org/zap"

	config "github.com/adzpm/telegram-clicker/internal/config"
	storage "github.com/adzpm/telegram-clicker/internal/storage"
)

type (
	REST struct {
		srv *fiber.App
		lgr *zap.Logger
		str *storage.Storage
		cfg *config.REST
	}
)

func New(lgr *zap.Logger, str *storage.Storage, cfg *config.REST) *REST {
	return &REST{
		srv: fiber.New(),
		lgr: lgr,
		cfg: cfg,
		str: str,
	}
}

func (a *REST) setupRoutes(ctx context.Context) {
	a.lgr.Debug("setting up routes")

	a.srv.Static("/", a.cfg.WebPath)

	a.srv.Get("/enter", a.EnterGame)
	a.srv.Get("/click", a.ClickCard)
	a.srv.Get("/buy", a.BuyCard)
	a.srv.Get("/reset", a.ResetGame)
}

func (a *REST) Start(ctx context.Context) error {
	a.setupRoutes(ctx)

	return a.srv.Listen(a.cfg.Host + ":" + a.cfg.Port)
}
