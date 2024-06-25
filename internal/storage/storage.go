package storage

import (
	"log"
	"os"
	"time"

	zap "go.uber.org/zap"
	sqlite "gorm.io/driver/sqlite"
	gorm "gorm.io/gorm"
	logger "gorm.io/gorm/logger"

	model "github.com/adzpm/telegram-clicker/internal/model"
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
		model.UserProduct{},
		model.Product{},
	}
)

func NewStorage(lgr *zap.Logger, cfg *Config) (*Storage, error) {
	str, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,          // Don't include params in the SQL log
			},
		),
	})
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

func (s *Storage) InsertUser(telegramID uint64) (*model.User, error) {
	res := s.str.Table("users").Create(&model.User{TelegramID: telegramID})
	if res.Error != nil {
		return nil, res.Error
	}

	user, err := s.SelectUser(telegramID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) SelectUser(telegramID uint64) (*model.User, error) {
	var user *model.User

	res := s.str.Table("users").Where("telegram_id = ?", telegramID).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (s *Storage) SelectUsers() ([]model.User, error) {
	var users []model.User

	res := s.str.Table("users").Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (s *Storage) UpdateUserCoins(telegramID, coins uint64) (*model.User, error) {
	res := s.str.Table("users").Where("telegram_id = ?", telegramID).Update("coins", coins)
	if res.Error != nil {
		return nil, res.Error
	}

	user, err := s.SelectUser(telegramID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) UpdateUserLastSeen(telegramID, lastSeen uint64) (*model.User, error) {
	res := s.str.Table("users").Where("telegram_id = ?", telegramID).Update("last_seen", lastSeen)
	if res.Error != nil {
		return nil, res.Error
	}

	user, err := s.SelectUser(telegramID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) InsertUserProduct(telegramID, productID, level uint64) (*model.UserProduct, error) {
	var userProduct *model.UserProduct

	res := s.str.Table("user_products").Create(&model.UserProduct{TelegramID: telegramID, ProductID: productID, Level: level})
	if res.Error != nil {
		return nil, res.Error
	}

	userProduct, err := s.SelectUserProduct(telegramID, productID)
	if err != nil {
		return nil, err
	}

	return userProduct, nil
}

func (s *Storage) SelectUserProduct(telegramID, productID uint64) (*model.UserProduct, error) {
	var userProduct *model.UserProduct

	res := s.str.Table("user_products").Where("telegram_id = ? AND product_id = ?", telegramID, productID).First(&userProduct)
	if res.Error != nil {
		return nil, res.Error
	}

	return userProduct, nil
}

func (s *Storage) SelectUserProducts(telegramID uint64) ([]model.UserProduct, error) {
	var userProducts []model.UserProduct

	res := s.str.Table("user_products").Where("telegram_id = ?", telegramID).Find(&userProducts)
	if res.Error != nil {
		return nil, res.Error
	}

	return userProducts, nil
}

func (s *Storage) UpdateUserProductLevel(telegramID, productID, level uint64) (*model.UserProduct, error) {
	res := s.str.Table("user_products").Where("telegram_id = ? AND product_id = ?", telegramID, productID).Update("level", level)
	if res.Error != nil {
		return nil, res.Error
	}

	userProduct, err := s.SelectUserProduct(telegramID, productID)
	if err != nil {
		return nil, err
	}

	return userProduct, nil
}

// InsertProduct - TODO: FIX THIS
func (s *Storage) InsertProduct(name, imageURL string, startPrice uint64, priceMultiplier float64, startCoins uint64, coinsMultiplier float64, maxLevel uint64) (*model.Product, error) {
	res := s.str.Table("products").Create(&model.Product{
		Name:            name,
		ImageURL:        imageURL,
		StartPrice:      startPrice,
		PriceMultiplier: priceMultiplier,
		StartCoins:      startCoins,
		CoinsMultiplier: coinsMultiplier,
		MaxLevel:        maxLevel,
	})

	if res.Error != nil {
		return nil, res.Error
	}

	product, err := s.SelectProduct(uint64(res.RowsAffected))
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Storage) SelectProduct(productID uint64) (*model.Product, error) {
	var products *model.Product

	res := s.str.Table("products").Where("id = ?", productID).First(&products)
	if res.Error != nil {
		return nil, res.Error
	}

	return products, nil
}

func (s *Storage) SelectProducts() ([]model.Product, error) {
	var products []model.Product

	res := s.str.Table("products").Find(&products)
	if res.Error != nil {
		return nil, res.Error
	}

	return products, nil
}
