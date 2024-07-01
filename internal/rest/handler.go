package rest

import (
	"errors"
	"net/http"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	zap "go.uber.org/zap"
	gorm "gorm.io/gorm"

	math "github.com/adzpm/telegram-clicker/internal/math"
	restModel "github.com/adzpm/telegram-clicker/internal/model/rest"
	storageModel "github.com/adzpm/telegram-clicker/internal/model/storage"
)

func mergeCards(
	user *storageModel.User,
	allCards []storageModel.Card,
	userCards []storageModel.UserCard,
) map[uint64]*restModel.GameCard {
	var (
		cards        = make(map[uint64]*restModel.GameCard, len(allCards))
		allCardsMap  = make(map[uint64]*storageModel.Card, len(allCards))
		userCardsMap = make(map[uint64]*storageModel.UserCard, len(userCards))
	)

	// prepare allCards map
	for _, card := range allCards {
		allCardsMap[card.ID] = &card
	}

	// prepare userCards map
	for _, userCard := range userCards {
		userCardsMap[userCard.CardID] = &userCard
	}

	// fill cards map
	for _, card := range allCards {
		cards[card.ID] = &rest.GameCard{
			ID:                     card.ID,
			Name:                   card.Name,
			ImageURL:               card.ImageURL,
			MaxLevel:               card.MaxLevel,
			NextLevelPrice:         card.Price,
			ClickTimeout:           card.ClickTimeout,
			NextLevelCoinsPerClick: card.CoinsPerClick,
		}
	}

	// update cards map with userCards data
	for _, userCard := range userCards {
		if _, ok := cards[userCard.CardID]; !ok {
			continue
		}

		if _, ok := allCardsMap[userCard.CardID]; !ok {
			continue
		}

		if userCard.Level < 1 {
			continue
		}

		var (
			cardID             = userCard.CardID
			level              = userCard.Level
			nextClick          = userCard.NextClick
			lastClick          = userCard.LastClick
			startPrice         = allCardsMap[cardID].Price
			startCoinsPerClick = allCardsMap[cardID].CoinsPerClick
			priceMp            = allCardsMap[cardID].PriceMultiplier
			invMp              = math.CalculateInvestorsMultiplier(user.Investors)
			nextPrice          = math.CalculateUpgradePrice(startPrice, level, priceMp)
			curPrice           = math.CalculateUpgradePrice(startPrice, level-1, priceMp)
			nextCoins          = math.CalculateAlgebraCoinsPerClick(startCoinsPerClick, level+1, invMp)
			curCoins           = math.CalculateAlgebraCoinsPerClick(startCoinsPerClick, level, invMp)
		)

		cards[cardID].CurrentLevel = level
		cards[cardID].CurrentPrice = curPrice
		cards[cardID].NextLevelPrice = nextPrice
		cards[cardID].CurrentCoinsPerClick = curCoins
		cards[cardID].NextLevelCoinsPerClick = nextCoins
		cards[cardID].NextClick = nextClick
		cards[cardID].LastClick = lastClick
	}

	return cards
}

func CreateGameResponse(user *model.User, allCards []model.Card, userCards []model.UserCard) *rest.Game {
	var (
		icount = math.CalculateInvestorsCount(user.EarnedCoins)
		curmlt = math.CalculateInvestorsMultiplier(user.Investors)
		nxtmlt = math.CalculateInvestorsMultiplier(icount)
	)

	return &rest.Game{
		UserID:                        user.ID,
		TelegramID:                    user.TelegramID,
		LastSeen:                      user.LastSeen,
		CurrentCoins:                  user.Coins,
		CurrentGold:                   user.Gold,
		CurrentInvestors:              user.Investors,
		CurrentInvestorsMultiplier:    curmlt,
		InvestorsMultiplierAfterReset: nxtmlt,
		InvestorsAfterReset:           icount,
		PercentsPerInvestor:           math.InvestorMultiplier * 100,
		Cards:                         mergeCards(user, allCards, userCards),
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

func (a *REST) EnterGame(c *fiber.Ctx) (err error) {
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

			if _, err = a.str.InsertUserCard(user.TelegramID, 1, 1); err != nil {
				return Throw500Error(c, err)
			}
		} else {
			return Throw500Error(c, err)
		}
	}

	var (
		timeNow   = uint64(time.Now().Unix())
		allCards  []model.Card
		userCards []model.UserCard
	)

	if allCards, err = a.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = a.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	// add coins for offline time

	if user, err = a.str.UpdateUserLastSeen(user.TelegramID, timeNow); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, CreateGameResponse(user, allCards, userCards))
}

func (a *REST) ClickCard(c *fiber.Ctx) (err error) {
	var (
		tn   = uint64(time.Now().Unix())
		tgID int
		cdID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, "telegram_id is required")
	}

	if cdID = c.QueryInt("card_id"); cdID == 0 {
		return Throw400Error(c, "card_id is required")
	}

	a.lgr.Info("try to click", zap.Int("telegram_id", tgID), zap.Int("card_id", cdID))

	var (
		user         *model.User
		card         *model.Card
		userCard     *model.UserCard
		coinsClicked uint64 = 0
	)

	if user, err = a.str.SelectUser(uint64(tgID)); err != nil {
		return Throw500Error(c, err)
	}

	if card, err = a.str.SelectCard(uint64(cdID)); err != nil {
		return Throw500Error(c, err)
	}

	if userCard, err = a.str.SelectUserCard(user.TelegramID, card.ID); err != nil {
		return Throw500Error(c, err)
	}

	if userCard.Level > 0 {
		coinsClicked = math.CalculateAlgebraCoinsPerClick(
			card.CoinsPerClick,
			userCard.Level,
			math.CalculateInvestorsMultiplier(user.Investors),
		)
	}

	if userCard.NextClick == 0 {
		userCard.NextClick = tn
	}

	if tn < userCard.NextClick {
		return Throw400Error(c, "you can't click now")
	}

	if user, err = a.str.UpdateUserCoins(user.TelegramID, user.Coins+coinsClicked); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = a.str.UpdateUserEarnedCoins(user.TelegramID, user.EarnedCoins+coinsClicked); err != nil {
		return Throw500Error(c, err)
	}

	if userCard, err = a.str.UpdateUserCardLastClick(user.TelegramID, card.ID, tn); err != nil {
		return Throw500Error(c, err)
	}

	if userCard, err = a.str.UpdateUserCardNextClick(user.TelegramID, card.ID, tn+card.ClickTimeout); err != nil {
		return Throw500Error(c, err)
	}

	var (
		allCards  []model.Card
		userCards []model.UserCard
	)

	if allCards, err = a.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = a.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, CreateGameResponse(user, allCards, userCards))
}

func (a *REST) BuyCard(c *fiber.Ctx) (err error) {
	var (
		tgID int
		prID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, "telegram_id is required")
	}

	if prID = c.QueryInt("card_id"); prID == 0 {
		return Throw400Error(c, "card_id is required")
	}

	a.lgr.Info("try to buy card", zap.Int("telegram_id", tgID), zap.Int("card_id", prID))

	var (
		user      *model.User
		card      *model.Card
		userCard  *model.UserCard
		allCards  []model.Card
		userCards []model.UserCard
	)

	if user, err = a.str.SelectUser(uint64(tgID)); err != nil {
		return Throw500Error(c, err)
	}

	if card, err = a.str.SelectCard(uint64(prID)); err != nil {
		return Throw500Error(c, err)
	}

	if allCards, err = a.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = a.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)

	}

	if userCard, err = a.str.SelectUserCard(user.TelegramID, card.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			priceToBuy := math.CalculateUpgradePrice(card.Price, 0, card.PriceMultiplier)

			if user.Coins < priceToBuy {
				return Throw400Error(c, "not enough coins")
			}

			if userCard, err = a.str.InsertUserCard(user.TelegramID, card.ID, 1); err != nil {
				return Throw500Error(c, err)
			}

			if user, err = a.str.UpdateUserCoins(user.TelegramID, user.Coins-priceToBuy); err != nil {
				return Throw500Error(c, err)
			}

			if userCards, err = a.str.SelectUserCards(user.TelegramID); err != nil {
				return Throw500Error(c, err)
			}

			return Throw200Response(c, CreateGameResponse(user, allCards, userCards))
		} else {
			return Throw500Error(c, err)
		}
	}

	priceToBuy := math.CalculateUpgradePrice(
		card.Price,
		userCard.Level,
		card.PriceMultiplier,
	)

	if user.Coins < priceToBuy {
		return Throw400Error(c, "not enough coins")
	}

	if user, err = a.str.UpdateUserCoins(user.TelegramID, user.Coins-priceToBuy); err != nil {
		return Throw500Error(c, err)
	}

	if userCard, err = a.str.UpdateUserCardLevel(user.TelegramID, card.ID, userCard.Level+1); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = a.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, CreateGameResponse(user, allCards, userCards))
}

func (a *REST) ResetGame(c *fiber.Ctx) (err error) {
	var (
		tgID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, "telegram_id is required")
	}

	a.lgr.Info("try to reset game", zap.Int("telegram_id", tgID))

	var (
		user      *model.User
		allCards  []model.Card
		userCards []model.UserCard
	)

	if user, err = a.str.SelectUser(uint64(tgID)); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = a.str.UpdateUserInvestors(user.TelegramID, math.CalculateInvestorsCount(user.EarnedCoins)); err != nil {
		return Throw500Error(c, err)
	}

	if allCards, err = a.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = a.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = a.str.UpdateUserEarnedCoins(user.TelegramID, 0); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = a.str.UpdateUserCoins(user.TelegramID, 0); err != nil {
		return Throw500Error(c, err)
	}

	for _, userCard := range userCards {
		if userCard.Level > 0 {
			if _, err = a.str.UpdateUserCardLevel(user.TelegramID, userCard.CardID, 0); err != nil {
				return Throw500Error(c, err)
			}
		}
	}

	if _, err = a.str.UpdateUserCardLevel(user.TelegramID, 1, 1); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = a.str.SelectUser(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	if allCards, err = a.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = a.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, CreateGameResponse(user, allCards, userCards))
}
