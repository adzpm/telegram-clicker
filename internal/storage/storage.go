package storage

import (
	model "github.com/adzpm/tg-clicker/internal/model"
	zap "go.uber.org/zap"
	sqlite "gorm.io/driver/sqlite"
	gorm "gorm.io/gorm"
	"time"
)

type (
	Config struct {
		Path string
	}

	Storage struct {
		str *gorm.DB
		lgr *zap.Logger
		cfg *Config
	}
)

var (
	migrate = []interface{}{
		model.User{},
	}
)

func NewStorage(lgr *zap.Logger, cfg *Config) (*Storage, error) {
	str, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = str.AutoMigrate(migrate...); err != nil {
		return nil, err
	}

	return &Storage{
		str: str,
		lgr: lgr,
		cfg: cfg,
	}, nil
}

func (s *Storage) CreateUser(u *model.User) (ru *model.User, err error) {
	res := s.str.Create(u)

	if res.Error != nil {
		s.lgr.Error("failed to create user", zap.Error(res.Error))
		return nil, res.Error
	}

	return s.GetUserByTelegramID(u.TelegramID)
}

func (s *Storage) GetUserByTelegramID(tid uint64) (u *model.User, err error) {
	res := s.str.Where("telegram_id = ?", tid).First(&u)

	if res.Error != nil {
		s.lgr.Error("failed to get user by id", zap.Error(res.Error))
		return nil, res.Error
	}

	return u, nil
}

func (s *Storage) AddCoinsByTelegramID(tid, cnt uint64) (u *model.User, err error) {
	u, err = s.GetUserByTelegramID(tid)
	if err != nil {
		return nil, err
	}

	u.Coins += cnt

	res := s.str.Save(u)

	if res.Error != nil {
		s.lgr.Error("failed to update", zap.Error(res.Error))
		return nil, res.Error
	}

	return u, nil
}

func (s *Storage) Login(tid uint64) (u *model.User, err error) {
	u, err = s.GetUserByTelegramID(tid)
	if err != nil {
		return nil, err
	}

	now := uint64(time.Now().Unix())

	// add 1 coin every 1min of inactivity. Max 3h
	if u.LastSeen == 0 {
		u.LastSeen = now
	} else {
		// 1 coin every 1min
		coins := (now - u.LastSeen) / 60
		if coins > 3*60 {
			coins = 3 * 60
		}

		u.Coins += coins
		u.LastSeen = now
	}

	res := s.str.Save(u)

	if res.Error != nil {
		s.lgr.Error("failed to update", zap.Error(res.Error))
		return nil, res.Error
	}

	return u, nil
}
