package api

import (
	"errors"
	"gorm.io/gorm"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	zap "go.uber.org/zap"

	"github.com/adzpm/telegram-clicker/internal/math"
	"github.com/adzpm/telegram-clicker/internal/model"
)

func mergeProducts(
	allProducts []model.Product,
	userProducts []model.UserProduct,
) map[uint64]*model.GameProduct {
	products := make(map[uint64]*model.GameProduct)

	for _, product := range allProducts {
		products[product.ID] = &model.GameProduct{
			ID:             product.ID,
			Name:           product.Name,
			ImageURL:       product.ImageURL,
			UpgradePrice:   math.CalculateUpgradePrice(product.StartPrice, 0, product.PriceMultiplier),
			CoinsPerClick:  math.CalculateCoinsPerClick(product.StartCoins, 0, product.CoinsMultiplier),
			CoinsPerMinute: 0, // TODO: IMPLEMENT THIS
			CurrentLevel:   0, // TODO: IMPLEMENT THIS
			MaxLevel:       product.MaxLevel,
		}
	}

	for _, userProduct := range userProducts {
		products[userProduct.ProductID].CurrentLevel = userProduct.Level

		products[userProduct.ProductID].UpgradePrice = math.CalculateUpgradePrice(
			allProducts[userProduct.ProductID].StartPrice,
			products[userProduct.ProductID].CurrentLevel,
			allProducts[userProduct.ProductID].PriceMultiplier,
		)

		products[userProduct.ProductID].CoinsPerClick = math.CalculateCoinsPerClick(
			allProducts[userProduct.ProductID].StartCoins,
			products[userProduct.ProductID].CurrentLevel,
			allProducts[userProduct.ProductID].CoinsMultiplier,
		)

		products[userProduct.ProductID].CoinsPerMinute = 0 // TODO: IMPLEMENT THIS
	}

	return products
}

func Throw500Error(c *fiber.Ctx, dst interface{}) (err error) {
	if err = c.Status(500).JSON(fiber.Map{"error": dst}); err != nil {
		return err
	}

	return nil
}

func Throw400Error(c *fiber.Ctx, dst interface{}) (err error) {
	if err = c.Status(400).JSON(fiber.Map{"error": dst}); err != nil {
		return err
	}

	return nil
}

func (a *API) Enter(c *fiber.Ctx) (err error) {
	var (
		user *model.User
		tgID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, "telegram_id is required")
	}

	a.lgr.Info("try to enter game", zap.Int("telegram_id", tgID))

	if user, err = a.str.SelectUser(uint64(tgID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.lgr.Warn("error while selecting user. Try to create new account", zap.Error(err))

			if user, err = a.str.InsertUser(uint64(tgID)); err != nil {
				return Throw500Error(c, err.Error())
			}

			if _, err = a.str.InsertUserProduct(user.TelegramID, 1, 1); err != nil {
				return Throw500Error(c, err.Error())
			}
		} else {
			return Throw500Error(c, err.Error())
		}
	}

	var (
		timeNow      = uint64(time.Now().Unix())
		allProducts  []model.Product
		userProducts []model.UserProduct
	)

	if allProducts, err = a.str.SelectProducts(); err != nil {
		return Throw500Error(c, err.Error())
	}

	if userProducts, err = a.str.SelectUserProducts(user.TelegramID); err != nil {
		return Throw500Error(c, err.Error())
	}

	// add coins for offline time

	if user, err = a.str.UpdateUserLastSeen(user.TelegramID, timeNow); err != nil {
		return Throw500Error(c, err.Error())
	}

	products := mergeProducts(allProducts, userProducts)

	response := &model.Game{
		UserID:       user.ID,
		TelegramID:   user.TelegramID,
		LastSeen:     user.LastSeen,
		CurrentCoins: user.Coins,
		CurrentGold:  0,
		Products:     products,
	}

	return c.Status(200).JSON(response)
}

func (a *API) Click(c *fiber.Ctx) (err error) {
	var (
		tgID int
		prID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, "telegram_id is required")
	}

	if prID = c.QueryInt("product_id"); prID == 0 {
		return Throw400Error(c, "product_id is required")
	}

	a.lgr.Info("try to click", zap.Int("telegram_id", tgID), zap.Int("product_id", prID))

	var (
		user         *model.User
		product      *model.Product
		userProduct  *model.UserProduct
		coinsClicked uint64 = 0
	)

	if user, err = a.str.SelectUser(uint64(tgID)); err != nil {
		return Throw500Error(c, err.Error())
	}

	if product, err = a.str.SelectProduct(uint64(prID)); err != nil {
		return Throw500Error(c, err.Error())
	}

	if userProduct, err = a.str.SelectUserProduct(user.TelegramID, product.ID); err != nil {
		return Throw500Error(c, err.Error())
	}

	if userProduct.Level > 0 {
		coinsClicked = math.CalculateCoinsPerClick(product.StartCoins, userProduct.Level, product.CoinsMultiplier)
	}

	if user, err = a.str.UpdateUserCoins(user.TelegramID, user.Coins+coinsClicked); err != nil {
		return Throw500Error(c, err.Error())
	}

	var (
		allProducts  []model.Product
		userProducts []model.UserProduct
	)

	if allProducts, err = a.str.SelectProducts(); err != nil {
		return Throw500Error(c, err.Error())
	}

	if userProducts, err = a.str.SelectUserProducts(user.TelegramID); err != nil {
		return Throw500Error(c, err.Error())
	}

	response := &model.Game{
		UserID:       user.ID,
		TelegramID:   user.TelegramID,
		LastSeen:     user.LastSeen,
		CurrentCoins: user.Coins,
		Products:     mergeProducts(allProducts, userProducts),
	}

	return c.Status(200).JSON(response)
}
