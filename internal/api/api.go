package api

import (
	"context"
	fiber "github.com/gofiber/fiber/v2"
	zap "go.uber.org/zap"

	storage "github.com/adzpm/telegram-clicker/internal/storage"
)

type (
	Config struct {
		Port    string
		WebPath string
	}

	API struct {
		srv *fiber.App
		lgr *zap.Logger
		str *storage.Storage
		cfg *Config
	}
)

func NewAPI(lgr *zap.Logger, str *storage.Storage, cfg *Config) *API {
	srv := fiber.New()

	return &API{
		srv: srv,
		lgr: lgr,
		cfg: cfg,
		str: str,
	}
}

func (a *API) setupRoutes(ctx context.Context) {
	a.lgr.Info("setting up routes")

	a.srv.Static("/", a.cfg.WebPath)

	a.srv.Get("/enter", a.Enter)
	a.srv.Get("/click", a.ClickProduct)
	a.srv.Get("/buy", a.BuyProduct)
}

func (a *API) Start(ctx context.Context) error {
	a.setupRoutes(ctx)

	return a.srv.Listen(":" + a.cfg.Port)
}
