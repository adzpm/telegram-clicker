package storage

import (
	zap "go.uber.org/zap"
	sqlite "gorm.io/driver/sqlite"
	gorm "gorm.io/gorm"

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

const (
	sqlInsertUser         = `INSERT INTO users (telegram_id) VALUES (?)`
	sqlSelectUser         = `SELECT * FROM users WHERE telegram_id = ?`
	sqlSelectUsers        = `SELECT * FROM users`
	sqlUpdateUserCoins    = `UPDATE users SET coins = ? WHERE telegram_id = ?`
	sqlUpdateUserLastSeen = `UPDATE users SET last_seen = ? WHERE telegram_id = ?`

	sqlInsertUserProduct      = `INSERT INTO user_products (user_id, product_id, lvl) VALUES (?, ?, ?)`
	sqlSelectUserProduct      = `SELECT * FROM user_products WHERE user_id = ? AND product_id = ?`
	sqlSelectUserProducts     = `SELECT * FROM user_products WHERE user_id = ?`
	sqlUpdateUserProductLevel = `UPDATE user_products SET level = ? WHERE user_id = ? AND product_id = ?`

	sqlInsertProduct  = `INSERT INTO products (name, image_url, start_price, price_multiplier, start_coins, coins_multiplier, max_level) VALUES (?, ?, ?, ?, ?, ?, ?)`
	sqlSelectProducts = `SELECT * FROM products`
	sqlSelectProduct  = `SELECT * FROM products WHERE id = ?`
)

func (s *Storage) InsertUser(telegramID uint64) (*model.User, error) {
	res := s.str.Exec(sqlInsertUser, telegramID)
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
	var user model.User

	res := s.str.Raw(sqlSelectUser, telegramID).Scan(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func (s *Storage) SelectUsers() ([]model.User, error) {
	var users []model.User

	res := s.str.Raw(sqlSelectUsers).Scan(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (s *Storage) UpdateUserCoins(telegramID, coins uint64) (*model.User, error) {
	res := s.str.Exec(sqlUpdateUserCoins, coins, telegramID)
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
	res := s.str.Exec(sqlUpdateUserLastSeen, lastSeen, telegramID)
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

	res := s.str.Exec(sqlInsertUserProduct, telegramID, productID, level)
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

	res := s.str.Raw(sqlSelectUserProduct, telegramID, productID).Scan(&userProduct)
	if res.Error != nil {
		return nil, res.Error
	}

	return userProduct, nil
}

func (s *Storage) SelectUserProducts(telegramID uint64) ([]model.UserProduct, error) {
	var userProducts []model.UserProduct

	res := s.str.Raw(sqlSelectUserProducts, telegramID).Scan(&userProducts)
	if res.Error != nil {
		return nil, res.Error
	}

	return userProducts, nil
}

func (s *Storage) UpdateUserProductLevel(telegramID, productID, level uint64) (*model.UserProduct, error) {
	res := s.str.Exec(sqlUpdateUserProductLevel, level, telegramID, productID)
	if res.Error != nil {
		return nil, res.Error
	}

	userProduct, err := s.SelectUserProduct(telegramID, productID)
	if err != nil {
		return nil, err
	}

	return userProduct, nil
}

func (s *Storage) InsertProduct(name, imageURL string, startPrice uint64, priceMultiplier float64, startCoins uint64, coinsMultiplier float64, maxLevel uint64) (*model.Product, error) {
	res := s.str.Exec(sqlInsertProduct, name, imageURL, startPrice, priceMultiplier, startCoins, coinsMultiplier, maxLevel)
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

	res := s.str.Raw(sqlSelectProduct, productID).Scan(&products)
	if res.Error != nil {
		return nil, res.Error
	}

	return products, nil
}

func (s *Storage) SelectProducts() ([]model.Product, error) {
	var products []model.Product

	res := s.str.Raw(sqlSelectProducts).Scan(&products)
	if res.Error != nil {
		return nil, res.Error
	}

	return products, nil
}
