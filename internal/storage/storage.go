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

	s.lgr.Info("filling products from file")

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

func (s *Storage) InsertUser(telegramID uint64) (*model.User, error) {
	s.lgr.Info("inserting user", zap.Uint64("telegram_id", telegramID))

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
	s.lgr.Info("selecting user", zap.Uint64("telegram_id", telegramID))

	var user *model.User

	res := s.str.Table("users").Where("telegram_id = ?", telegramID).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (s *Storage) SelectUsers() ([]model.User, error) {
	s.lgr.Info("selecting all users")

	var users []model.User

	res := s.str.Table("users").Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (s *Storage) UpdateUserCoins(telegramID, coins uint64) (*model.User, error) {
	s.lgr.Info("updating user coins",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("coins", coins),
	)

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
	s.lgr.Info("updating user last seen", zap.Uint64("telegram_id", telegramID), zap.Uint64("last_seen", lastSeen))

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
	s.lgr.Info("inserting user product",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("product_id", productID),
		zap.Uint64("level", level),
	)

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
	s.lgr.Info("selecting user product",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("product_id", productID),
	)

	var userProduct *model.UserProduct

	res := s.str.Table("user_products").Where("telegram_id = ? AND product_id = ?", telegramID, productID).First(&userProduct)
	if res.Error != nil {
		return nil, res.Error
	}

	return userProduct, nil
}

func (s *Storage) SelectUserProducts(telegramID uint64) ([]model.UserProduct, error) {
	s.lgr.Info("selecting user products", zap.Uint64("telegram_id", telegramID))

	var userProducts []model.UserProduct

	res := s.str.Table("user_products").Where("telegram_id = ?", telegramID).Find(&userProducts)
	if res.Error != nil {
		return nil, res.Error
	}

	return userProducts, nil
}

func (s *Storage) UpdateUserProductLevel(telegramID, productID, level uint64) (*model.UserProduct, error) {
	s.lgr.Info("updating user product level",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("product_id", productID),
		zap.Uint64("level", level),
	)

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
func (s *Storage) InsertProduct(name, imageURL string,
	startPrice uint64,
	priceMultiplier float64,
	startCoinsPerClick uint64,
	coinsMultiplier float64,
	maxLevel uint64,
) (*model.Product, error) {
	s.lgr.Info("inserting product",
		zap.String("name", name),
		zap.String("image_url", imageURL),
		zap.Uint64("start_price", startPrice),
		zap.Float64("price_multiplier", priceMultiplier),
		zap.Uint64("start_coins_per_click", startCoinsPerClick),
		zap.Float64("coins_multiplier", coinsMultiplier),
		zap.Uint64("max_level", maxLevel),
	)

	res := s.str.Table("products").Create(&model.Product{
		Name:                    name,
		ImageURL:                imageURL,
		StartProductPrice:       startPrice,
		ProductPriceMultiplier:  priceMultiplier,
		StartCoinsPerClick:      startCoinsPerClick,
		CoinsPerClickMultiplier: coinsMultiplier,
		MaxLevel:                maxLevel,
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
	s.lgr.Info("selecting product", zap.Uint64("product_id", productID))

	var products *model.Product

	res := s.str.Table("products").Where("id = ?", productID).First(&products)
	if res.Error != nil {
		return nil, res.Error
	}

	return products, nil
}

func (s *Storage) SelectProducts() ([]model.Product, error) {
	s.lgr.Info("selecting all products")

	var products []model.Product

	res := s.str.Table("products").Find(&products)
	if res.Error != nil {
		return nil, res.Error
	}

	return products, nil
}
