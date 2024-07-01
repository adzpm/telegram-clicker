package rest

import (
	"errors"
	"net/http"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	zap "go.uber.org/zap"
	gorm "gorm.io/gorm"

	restModel "github.com/adzpm/telegram-clicker/internal/model/rest"
	storageModel "github.com/adzpm/telegram-clicker/internal/model/storage"
)

const (
	keyError = "error"

	ErrorNotEnoughCoins       = "not enough coins"
	ErrorCantClickNow         = "you can't click now"
	ErrorTelegramIDIsRequired = "telegram_id is required"
	ErrorCardIDIsRequired     = "card_id is required"
)

func (r *REST) mergeCards(
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
		cards[card.ID] = &restModel.GameCard{
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
			invMp              = r.mth.CalculateInvestorsMultiplier(user.Investors)
			nextPrice          = r.mth.CalculateUpgradePrice(startPrice, level, priceMp)
			curPrice           = r.mth.CalculateUpgradePrice(startPrice, level-1, priceMp)
			nextCoins          = r.mth.CalculateAlgebraCoinsPerClick(startCoinsPerClick, level+1, invMp)
			curCoins           = r.mth.CalculateAlgebraCoinsPerClick(startCoinsPerClick, level, invMp)
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

func (r *REST) createGameResponse(user *storageModel.User, allCards []storageModel.Card, userCards []storageModel.UserCard) *restModel.Game {
	var (
		icount = r.mth.CalculateInvestorsCount(user.EarnedCoins)
		curmlt = r.mth.CalculateInvestorsMultiplier(user.Investors)
		nxtmlt = r.mth.CalculateInvestorsMultiplier(icount)
	)

	return &restModel.Game{
		UserID:                        user.ID,
		TelegramID:                    user.TelegramID,
		LastSeen:                      user.LastSeen,
		CurrentCoins:                  user.Coins,
		CurrentGold:                   user.Gold,
		CurrentInvestors:              user.Investors,
		CurrentInvestorsMultiplier:    curmlt,
		InvestorsMultiplierAfterReset: nxtmlt,
		InvestorsAfterReset:           icount,
		PercentsPerInvestor:           uint64(r.mth.GetGameVariables().PercentsForInvestor * 100),
		Cards:                         r.mergeCards(user, allCards, userCards),
	}
}

func Throw500Error(c *fiber.Ctx, dst interface{}) (err error) {
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{keyError: dst})
}

func Throw400Error(c *fiber.Ctx, dst interface{}) (err error) {
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{keyError: dst})
}

func Throw200Response(c *fiber.Ctx, dst interface{}) (err error) {
	return c.Status(http.StatusOK).JSON(dst)
}

func (r *REST) EnterGame(c *fiber.Ctx) (err error) {
	var (
		user *storageModel.User
		tgID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, ErrorTelegramIDIsRequired)
	}

	r.lgr.Info("try to enter game", zap.Int("telegram_id", tgID))

	if user, err = r.str.SelectUser(uint64(tgID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.lgr.Warn("error while selecting user. Try to create new account", zap.Error(err))

			if user, err = r.str.InsertUser(uint64(tgID), 0, 1000, 0); err != nil {
				return Throw500Error(c, err)
			}

			if _, err = r.str.InsertUserCard(user.TelegramID, 1, 1); err != nil {
				return Throw500Error(c, err)
			}
		} else {
			return Throw500Error(c, err)
		}
	}

	var (
		timeNow   = uint64(time.Now().Unix())
		allCards  []storageModel.Card
		userCards []storageModel.UserCard
	)

	if allCards, err = r.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = r.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	// add coins for offline time

	if user, err = r.str.UpdateUserLastSeen(user.TelegramID, timeNow); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, r.createGameResponse(user, allCards, userCards))
}

func (r *REST) ClickCard(c *fiber.Ctx) (err error) {
	var (
		tn   = uint64(time.Now().Unix())
		tgID int
		cdID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, ErrorTelegramIDIsRequired)
	}

	if cdID = c.QueryInt("card_id"); cdID == 0 {
		return Throw400Error(c, ErrorCardIDIsRequired)
	}

	r.lgr.Info("try to click", zap.Int("telegram_id", tgID), zap.Int("card_id", cdID))

	var (
		user         *storageModel.User
		card         *storageModel.Card
		userCard     *storageModel.UserCard
		coinsClicked uint64 = 0
	)

	if user, err = r.str.SelectUser(uint64(tgID)); err != nil {
		return Throw500Error(c, err)
	}

	if card, err = r.str.SelectCard(uint64(cdID)); err != nil {
		return Throw500Error(c, err)
	}

	if userCard, err = r.str.SelectUserCard(user.TelegramID, card.ID); err != nil {
		return Throw500Error(c, err)
	}

	if userCard.Level > 0 {
		coinsClicked = r.mth.CalculateAlgebraCoinsPerClick(
			card.CoinsPerClick,
			userCard.Level,
			r.mth.CalculateInvestorsMultiplier(user.Investors),
		)
	}

	if userCard.NextClick == 0 {
		userCard.NextClick = tn
	}

	if tn < userCard.NextClick {
		return Throw400Error(c, ErrorCantClickNow)
	}

	if user, err = r.str.UpdateUserCoins(user.TelegramID, user.Coins+coinsClicked); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = r.str.UpdateUserEarnedCoins(user.TelegramID, user.EarnedCoins+coinsClicked); err != nil {
		return Throw500Error(c, err)
	}

	if userCard, err = r.str.UpdateUserCardLastClick(user.TelegramID, card.ID, tn); err != nil {
		return Throw500Error(c, err)
	}

	if userCard, err = r.str.UpdateUserCardNextClick(user.TelegramID, card.ID, tn+card.ClickTimeout); err != nil {
		return Throw500Error(c, err)
	}

	var (
		allCards  []storageModel.Card
		userCards []storageModel.UserCard
	)

	if allCards, err = r.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = r.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, r.createGameResponse(user, allCards, userCards))
}

func (r *REST) BuyCard(c *fiber.Ctx) (err error) {
	var (
		tgID int
		prID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, ErrorTelegramIDIsRequired)
	}

	if prID = c.QueryInt("card_id"); prID == 0 {
		return Throw400Error(c, ErrorCardIDIsRequired)
	}

	r.lgr.Info("try to buy card", zap.Int("telegram_id", tgID), zap.Int("card_id", prID))

	var (
		user      *storageModel.User
		card      *storageModel.Card
		userCard  *storageModel.UserCard
		allCards  []storageModel.Card
		userCards []storageModel.UserCard
	)

	if user, err = r.str.SelectUser(uint64(tgID)); err != nil {
		return Throw500Error(c, err)
	}

	if card, err = r.str.SelectCard(uint64(prID)); err != nil {
		return Throw500Error(c, err)
	}

	if allCards, err = r.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = r.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)

	}

	if userCard, err = r.str.SelectUserCard(user.TelegramID, card.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			priceToBuy := r.mth.CalculateUpgradePrice(card.Price, 0, card.PriceMultiplier)

			if user.Coins < priceToBuy {
				return Throw400Error(c, ErrorNotEnoughCoins)
			}

			if userCard, err = r.str.InsertUserCard(user.TelegramID, card.ID, 1); err != nil {
				return Throw500Error(c, err)
			}

			if user, err = r.str.UpdateUserCoins(user.TelegramID, user.Coins-priceToBuy); err != nil {
				return Throw500Error(c, err)
			}

			if userCards, err = r.str.SelectUserCards(user.TelegramID); err != nil {
				return Throw500Error(c, err)
			}

			return Throw200Response(c, r.createGameResponse(user, allCards, userCards))
		} else {
			return Throw500Error(c, err)
		}
	}

	priceToBuy := r.mth.CalculateUpgradePrice(
		card.Price,
		userCard.Level,
		card.PriceMultiplier,
	)

	if user.Coins < priceToBuy {
		return Throw400Error(c, ErrorNotEnoughCoins)
	}

	if user, err = r.str.UpdateUserCoins(user.TelegramID, user.Coins-priceToBuy); err != nil {
		return Throw500Error(c, err)
	}

	if userCard, err = r.str.UpdateUserCardLevel(user.TelegramID, card.ID, userCard.Level+1); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = r.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, r.createGameResponse(user, allCards, userCards))
}

func (r *REST) ResetGame(c *fiber.Ctx) (err error) {
	var (
		tgID int
	)

	if tgID = c.QueryInt("telegram_id"); tgID == 0 {
		return Throw400Error(c, ErrorTelegramIDIsRequired)
	}

	r.lgr.Info("try to reset game", zap.Int("telegram_id", tgID))

	var (
		user      *storageModel.User
		allCards  []storageModel.Card
		userCards []storageModel.UserCard
	)

	if user, err = r.str.SelectUser(uint64(tgID)); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = r.str.UpdateUserInvestors(user.TelegramID, r.mth.CalculateInvestorsCount(user.EarnedCoins)); err != nil {
		return Throw500Error(c, err)
	}

	if allCards, err = r.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = r.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = r.str.UpdateUserEarnedCoins(user.TelegramID, 0); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = r.str.UpdateUserCoins(user.TelegramID, 0); err != nil {
		return Throw500Error(c, err)
	}

	for _, userCard := range userCards {
		if userCard.Level > 0 {
			if _, err = r.str.UpdateUserCardLevel(user.TelegramID, userCard.CardID, 0); err != nil {
				return Throw500Error(c, err)
			}

			if _, err = r.str.UpdateUserCardLastClick(user.TelegramID, userCard.CardID, 0); err != nil {
				return Throw500Error(c, err)
			}

			if _, err = r.str.UpdateUserCardNextClick(user.TelegramID, userCard.CardID, 0); err != nil {
				return Throw500Error(c, err)
			}
		}
	}

	if _, err = r.str.UpdateUserCardLevel(user.TelegramID, 1, 1); err != nil {
		return Throw500Error(c, err)
	}

	if user, err = r.str.SelectUser(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	if allCards, err = r.str.SelectCards(); err != nil {
		return Throw500Error(c, err)
	}

	if userCards, err = r.str.SelectUserCards(user.TelegramID); err != nil {
		return Throw500Error(c, err)
	}

	return Throw200Response(c, r.createGameResponse(user, allCards, userCards))
}
