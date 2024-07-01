package rest

import (
	"context"
	"github.com/adzpm/telegram-clicker/internal/math"

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
		mth *math.Math
	}
)

func New(lgr *zap.Logger, str *storage.Storage, mth *math.Math, cfg *config.REST) *REST {
	return &REST{
		srv: fiber.New(),
		lgr: lgr,
		cfg: cfg,
		mth: mth,
		str: str,
	}
}

func (r *REST) setupRoutes(ctx context.Context) {
	r.lgr.Debug("setting up routes")

	r.srv.Static("/", r.cfg.WebPath)

	r.srv.Get("/enter", r.EnterGame)
	r.srv.Get("/click", r.ClickCard)
	r.srv.Get("/buy", r.BuyCard)
	r.srv.Get("/reset", r.ResetGame)
}

func (r *REST) Start(ctx context.Context) error {
	r.setupRoutes(ctx)

	return r.srv.Listen(r.cfg.Host + ":" + r.cfg.Port)
}
