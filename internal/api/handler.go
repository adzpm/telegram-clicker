package api

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
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
	var (
		products        = make(map[uint64]*model.GameProduct, len(allProducts))
		allProductsMap  = make(map[uint64]*model.Product, len(allProducts))
		userProductsMap = make(map[uint64]*model.UserProduct, len(userProducts))
	)

	// prepare allProducts map
	for _, product := range allProducts {
		allProductsMap[product.ID] = &product
	}

	// prepare userProducts map
	for _, userProduct := range userProducts {
		userProductsMap[userProduct.ProductID] = &userProduct
	}

	// fill products map
	for _, product := range allProducts {
		products[product.ID] = &model.GameProduct{
			ID:                     product.ID,
			Name:                   product.Name,
			ImageURL:               product.ImageURL,
			CurrentLevel:           0,
			MaxLevel:               product.MaxLevel,
			CurrentPrice:           0,
			NextLevelPrice:         product.StartProductPrice,
			CurrentCoinsPerClick:   0,
			NextLevelCoinsPerClick: product.StartCoinsPerClick,
		}
	}

	// update products map with userProducts data
	for _, userProduct := range userProducts {
		var (
			productID          = userProduct.ProductID
			level              = userProduct.Level
			startPrice         = allProductsMap[productID].StartProductPrice
			startCoinsPerClick = allProductsMap[productID].StartCoinsPerClick
			priceMp            = allProductsMap[productID].ProductPriceMultiplier
			coinsMP            = allProductsMap[productID].CoinsPerClickMultiplier
			nextPrice          = math.CalculateUpgradePrice(startPrice, level, priceMp)
			curPrice           = math.CalculateUpgradePrice(startPrice, level-1, priceMp)
			nextCoins          = math.CalculateCoinsPerClick(startCoinsPerClick, level+1, coinsMP)
			curCoins           = math.CalculateCoinsPerClick(startCoinsPerClick, level, coinsMP)
		)

		products[productID].CurrentLevel = level
		products[productID].CurrentPrice = curPrice
		products[productID].NextLevelPrice = nextPrice
		products[productID].CurrentCoinsPerClick = curCoins
		products[productID].NextLevelCoinsPerClick = nextCoins
	}

	return products
}

func CreateGameResponse(user *model.User, allProducts []model.Product, userProducts []model.UserProduct) *model.Game {
	return &model.Game{
		UserID:           user.ID,
		TelegramID:       user.TelegramID,
		LastSeen:         user.LastSeen,
		CurrentCoins:     user.Coins,
		CurrentGold:      user.Gold,
		CurrentInvestors: user.Investors,
		Products:         mergeProducts(allProducts, userProducts),
	}
}

func Throw500Error(c *fiber.Ctx, dst interface{}) (err error) {
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": dst})
}

func Throw400Error(c *fiber.Ctx, dst interface{}) (err error) {
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": dst})
}

func Throw200Response(c *fiber.Ctx, dst interface{}) (err error) {
	return c.Status(http.StatusOK).JSON(dst)
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

			if user, err = a.str.InsertUser(uint64(tgID), 0, 1000, 0); err != nil {
				return Throw500Error(c, err)
			}

			if _, err = a.str.InsertUserProduct(user.TelegramID, 1, 1); err != nil {
				return Throw500Error(c, err)
			}
		} else {
			return Throw500Error(c, err)
		}
	}

	var (
		timeNow      = uint64(time.Now().Unix())
		allProducts  []model.Product
		userProducts []model.UserProduct
	)

	if allProducts, err = a.str.SelectProducts(); err != nil {
		return Throw500Error(c, err)
	}

	if userProducts, err = a.str.SelectUserProducts(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	// add coins for offline time

	if user, err = a.str.UpdateUserLastSeen(user.TelegramID, timeNow); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, CreateGameResponse(user, allProducts, userProducts))
}

func (a *API) ClickProduct(c *fiber.Ctx) (err error) {
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
		return Throw500Error(c, err)
	}

	if product, err = a.str.SelectProduct(uint64(prID)); err != nil {
		return Throw500Error(c, err)
	}

	if userProduct, err = a.str.SelectUserProduct(user.TelegramID, product.ID); err != nil {
		return Throw500Error(c, err)
	}

	if userProduct.Level > 0 {
		coinsClicked = math.CalculateCoinsPerClick(
			product.StartCoinsPerClick,
			userProduct.Level,
			product.CoinsPerClickMultiplier,
		)
	}

	if user, err = a.str.UpdateUserCoins(user.TelegramID, user.Coins+coinsClicked); err != nil {
		return Throw500Error(c, err)
	}

	var (
		allProducts  []model.Product
		userProducts []model.UserProduct
	)

	if allProducts, err = a.str.SelectProducts(); err != nil {
		return Throw500Error(c, err)
	}

	if userProducts, err = a.str.SelectUserProducts(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, CreateGameResponse(user, allProducts, userProducts))
}

func (a *API) BuyProduct(c *fiber.Ctx) (err error) {
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

	a.lgr.Info("try to buy product", zap.Int("telegram_id", tgID), zap.Int("product_id", prID))

	var (
		user         *model.User
		product      *model.Product
		userProduct  *model.UserProduct
		allProducts  []model.Product
		userProducts []model.UserProduct
	)

	if user, err = a.str.SelectUser(uint64(tgID)); err != nil {
		return Throw500Error(c, err)
	}

	if product, err = a.str.SelectProduct(uint64(prID)); err != nil {
		return Throw500Error(c, err)
	}

	if allProducts, err = a.str.SelectProducts(); err != nil {
		return Throw500Error(c, err)
	}

	if userProducts, err = a.str.SelectUserProducts(user.TelegramID); err != nil {
		return Throw500Error(c, err)

	}

	if userProduct, err = a.str.SelectUserProduct(user.TelegramID, product.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			priceToBuy := math.CalculateUpgradePrice(product.StartProductPrice, 0, product.ProductPriceMultiplier)

			if user.Coins < priceToBuy {
				return Throw400Error(c, "not enough coins")
			}

			if userProduct, err = a.str.InsertUserProduct(user.TelegramID, product.ID, 1); err != nil {
				return Throw500Error(c, err)
			}

			if user, err = a.str.UpdateUserCoins(user.TelegramID, user.Coins-priceToBuy); err != nil {
				return Throw500Error(c, err)
			}

			if userProducts, err = a.str.SelectUserProducts(user.TelegramID); err != nil {
				return Throw500Error(c, err)
			}

			return Throw200Response(c, CreateGameResponse(user, allProducts, userProducts))
		} else {
			return Throw500Error(c, err)
		}
	}

	priceToBuy := math.CalculateUpgradePrice(
		product.StartProductPrice,
		userProduct.Level,
		product.ProductPriceMultiplier,
	)

	if user.Coins < priceToBuy {
		return Throw400Error(c, "not enough coins")
	}

	if user, err = a.str.UpdateUserCoins(user.TelegramID, user.Coins-priceToBuy); err != nil {
		return Throw500Error(c, err)
	}

	if userProduct, err = a.str.UpdateUserProductLevel(user.TelegramID, product.ID, userProduct.Level+1); err != nil {
		return Throw500Error(c, err)
	}

	if userProducts, err = a.str.SelectUserProducts(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, CreateGameResponse(user, allProducts, userProducts))
}
