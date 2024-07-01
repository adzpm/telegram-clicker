package main

import (
	"context"
	"os"

	zap "go.uber.org/zap"

	config "github.com/adzpm/telegram-clicker/internal/config"
	math "github.com/adzpm/telegram-clicker/internal/math"
	rest "github.com/adzpm/telegram-clicker/internal/rest"
	storage "github.com/adzpm/telegram-clicker/internal/storage"
)

const (
	envClickerConfigPath = "CLICKER_CONFIG_PATH"
	defClickerConfigPath = "example.config.yaml"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		cfgPath     = getEnv(envClickerConfigPath, defClickerConfigPath)
		cfg         = config.New()

		lgr *zap.Logger
		str *storage.Storage
		rst *rest.REST
		mth *math.Math
		err error
	)

	defer cancel()

	if err = cfg.Read(cfgPath); err != nil {
		panic(err)
	}

	if lgr, err = zap.NewProduction(); err != nil {
		panic(err)
	}

	defer func() { _ = lgr.Sync() }()

	if str, err = storage.New(lgr, &cfg.Storage); err != nil {
		panic(err)
	}

	if mth = math.New(&cfg.GameVariables); err != nil {
		panic(err)
	}

	rst = rest.New(lgr, str, mth, &cfg.REST)

	if err = rst.Start(ctx); err != nil {
		panic(err)
	}
}
