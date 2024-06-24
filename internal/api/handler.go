package api

import (
	"github.com/adzpm/telegram-clicker/internal/math"
	"github.com/adzpm/telegram-clicker/internal/model"
	fiber "github.com/gofiber/fiber/v2"
	zap "go.uber.org/zap"
	"time"
)

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
		a.lgr.Warn("error while selecting user. Try to create new account", zap.Error(err))

		if user, err = a.str.InsertUser(uint64(tgID)); err != nil {
			return Throw500Error(c, err.Error())
		}

		if _, err = a.str.InsertUserProduct(user.ID, 1, 1); err != nil {
			return Throw500Error(c, err.Error())
		}
	}

	var (
		timeNow        = uint64(time.Now().Unix())
		offlineMinutes = timeNow - user.LastSeen/60
		allProducts    []model.Product
		userProducts   []model.UserProduct
		coinsToAdd     uint64
	)

	// calculate coins for offline time
	if user.LastSeen != 0 {
		// if user was offline more than 60 minutes, we will calculate only 60 minutes
		if offlineMinutes > 60 {
			offlineMinutes = 60
		}

		if allProducts, err = a.str.SelectProducts(); err != nil {
			return Throw500Error(c, err.Error())
		}

		if userProducts, err = a.str.SelectUserProducts(user.ID); err != nil {
			return Throw500Error(c, err.Error())
		}

		for _, product := range allProducts {
			for _, userProduct := range userProducts {
				if product.ID == userProduct.ProductID {
					coinsToAdd += math.CalculateCoinsPerMinute(
						product.StartCoins,
						userProduct.Level,
						product.CoinsMultiplier,
					) * offlineMinutes
				}
			}
		}

		finalCoins := user.Coins + coinsToAdd

		if user, err = a.str.UpdateUserCoins(user.TelegramID, finalCoins); err != nil {
			return Throw500Error(c, err.Error())
		}
	}

	if user, err = a.str.UpdateUserLastSeen(user.TelegramID, timeNow); err != nil {
		return Throw500Error(c, err.Error())
	}

	products := make(map[uint64]*model.GameProduct)

	// fill products map
	for _, product := range allProducts {
		products[product.ID] = &model.GameProduct{
			ID:             product.ID,
			Name:           product.Name,
			ImageURL:       product.ImageURL,
			UpgradePrice:   math.CalculateUpgradePrice(product.StartPrice, 0, product.PriceMultiplier),
			CoinsPerClick:  math.CalculateCoinsPerClick(product.StartCoins, 0, product.CoinsMultiplier),
			CoinsPerMinute: math.CalculateCoinsPerMinute(product.StartCoins, 0, product.CoinsMultiplier),
			CurrentLevel:   0,
			MaxLevel:       product.MaxLevel,
		}
	}

	// fill current level, upgrade price, coins per click and coins per minute
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

		products[userProduct.ProductID].CoinsPerMinute = math.CalculateCoinsPerMinute(
			allProducts[userProduct.ProductID].StartCoins,
			products[userProduct.ProductID].CurrentLevel,
			allProducts[userProduct.ProductID].CoinsMultiplier,
		)
	}

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
