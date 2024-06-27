package storage

import (
	"encoding/json"
	"io"
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

	res := &Storage{
		str: str,
		lgr: lgr,
		cfg: cfg,
	}

	return res, res.FillProductsFromFileIfTableEmpty()
}

func (s *Storage) FillProductsFromFileIfTableEmpty() error {
	var products []model.Product

	res := s.str.Table("products").Find(&products)
	if res.Error != nil {
		return res.Error
	}

	if len(products) > 0 {
		return nil
	}

	s.lgr.Debug("filling products from file")

	file, err := os.Open("products.json")
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	fileBytes, err := io.ReadAll(file)

	if err = json.Unmarshal(fileBytes, &products); err != nil {
		return err
	}

	for _, product := range products {
		if res = s.str.Table("products").Create(&product); res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func (s *Storage) InsertUser(telegramID, coins, gold, investors uint64) (user *model.User, err error) {
	s.lgr.Debug("inserting user", zap.Uint64("telegram_id", telegramID))

	if res := s.str.Table("users").Create(&model.User{
		TelegramID: telegramID,
		LastSeen:   uint64(time.Now().Unix()),
		Coins:      coins,
		Gold:       gold,
		Investors:  investors,
	}); res.Error != nil {
		return nil, res.Error
	}

	if user, err = s.SelectUser(telegramID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) SelectUser(telegramID uint64) (user *model.User, err error) {
	s.lgr.Debug("selecting user", zap.Uint64("telegram_id", telegramID))

	res := s.str.Table("users").Where("telegram_id = ?", telegramID).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (s *Storage) SelectUsers() (users []model.User, err error) {
	s.lgr.Debug("selecting all users")

	if res := s.str.Table("users").Find(&users); res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (s *Storage) UpdateUserCoins(telegramID, coins uint64) (user *model.User, err error) {
	s.lgr.Debug("updating user coins",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("coins", coins),
	)

	if res := s.str.Table("users").Where("telegram_id = ?", telegramID).Update("coins", coins); res.Error != nil {
		return nil, res.Error
	}

	if user, err = s.SelectUser(telegramID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) UpdateUserLastSeen(telegramID, lastSeen uint64) (user *model.User, err error) {
	s.lgr.Debug("updating user last seen",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("last_seen", lastSeen),
	)

	if res := s.str.Table("users").Where("telegram_id = ?", telegramID).Update("last_seen", lastSeen); res.Error != nil {
		return nil, res.Error
	}

	if user, err = s.SelectUser(telegramID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) InsertUserProduct(telegramID, productID, level uint64) (userProduct *model.UserProduct, err error) {
	s.lgr.Debug("inserting user product",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("product_id", productID),
		zap.Uint64("level", level),
	)

	if res := s.str.Table("user_products").Create(&model.UserProduct{
		TelegramID: telegramID,
		ProductID:  productID,
		Level:      level,
	}); res.Error != nil {
		return nil, res.Error
	}

	if userProduct, err = s.SelectUserProduct(telegramID, productID); err != nil {
		return nil, err
	}

	return userProduct, nil
}

func (s *Storage) SelectUserProduct(telegramID, productID uint64) (userProduct *model.UserProduct, err error) {
	s.lgr.Debug("selecting user product",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("product_id", productID),
	)

	if res := s.str.Table("user_products").Where("telegram_id = ? AND product_id = ?", telegramID, productID).First(&userProduct); res.Error != nil {
		return nil, res.Error
	}

	return userProduct, nil
}

func (s *Storage) SelectUserProducts(telegramID uint64) (userProducts []model.UserProduct, err error) {
	s.lgr.Debug("selecting user products", zap.Uint64("telegram_id", telegramID))

	if res := s.str.Table("user_products").Where("telegram_id = ?", telegramID).Find(&userProducts); res.Error != nil {
		return nil, res.Error
	}

	return userProducts, nil
}

func (s *Storage) UpdateUserProductLevel(telegramID, productID, level uint64) (userProduct *model.UserProduct, err error) {
	s.lgr.Debug("updating user product level",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("product_id", productID),
		zap.Uint64("level", level),
	)

	if res := s.str.Table("user_products").Where("telegram_id = ? AND product_id = ?", telegramID, productID).Update("level", level); res.Error != nil {
		return nil, res.Error
	}

	if userProduct, err = s.SelectUserProduct(telegramID, productID); err != nil {
		return nil, err
	}

	return userProduct, nil
}

func (s *Storage) SelectProduct(productID uint64) (products *model.Product, err error) {
	s.lgr.Debug("selecting product", zap.Uint64("product_id", productID))

	if res := s.str.Table("products").Where("id = ?", productID).First(&products); res.Error != nil {
		return nil, res.Error
	}

	return products, nil
}

func (s *Storage) SelectProducts() (products []model.Product, err error) {
	s.lgr.Debug("selecting all products")

	if res := s.str.Table("products").Find(&products); res.Error != nil {
		return nil, res.Error
	}

	return products, nil
}
