package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	zap "go.uber.org/zap"
	postgres "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
	logger "gorm.io/gorm/logger"

	config "github.com/adzpm/telegram-clicker/internal/config"
	storageModel "github.com/adzpm/telegram-clicker/internal/model/storage"
)

type (
	Storage struct {
		str *gorm.DB
		lgr *zap.Logger
		cfg *config.Storage
	}
)

var (
	migrate = []interface{}{
		storageModel.User{},
		storageModel.UserCard{},
		storageModel.Card{},
	}
)

func New(lgr *zap.Logger, cfg *config.Storage) (_ *Storage, err error) {
	var (
		str *gorm.DB

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			cfg.Host,
			cfg.DBUser,
			cfg.DBPass,
			cfg.DBName,
			cfg.Port,
		)
	)

	if str, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		}),
	}); err != nil {
		return nil, err
	}

	if err = str.AutoMigrate(migrate...); err != nil {
		return nil, err
	}

	res := &Storage{
		str: str,
		lgr: lgr,
		cfg: cfg,
	}

	return res, res.FillCardsFromFileIfTableEmpty()
}

func (s *Storage) FillCardsFromFileIfTableEmpty() (err error) {
	var (
		cards []storageModel.Card
		fb    []byte
		res   *gorm.DB
		file  *os.File
	)

	if res = s.str.Table("cards").Find(&cards); res.Error != nil {
		return res.Error
	}

	if len(cards) > 0 {
		return nil
	}

	s.lgr.Debug("filling cards")

	if file, err = os.Open("cards.json"); err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if fb, err = io.ReadAll(file); err != nil {
		return err
	}

	if err = json.Unmarshal(fb, &cards); err != nil {
		return err
	}

	for _, card := range cards {
		if res = s.str.Table("cards").Create(&card); res.Error != nil {
			return res.Error
		}
	}

	return nil
}
