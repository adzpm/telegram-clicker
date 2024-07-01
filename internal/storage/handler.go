package storage

import (
	"github.com/adzpm/telegram-clicker/internal/model/storage"
	"go.uber.org/zap"
	"time"
)

func (s *Storage) InsertUser(telegramID, coins, gold, investors uint64) (user *storage.User, err error) {
	s.lgr.Debug("inserting user", zap.Uint64("telegram_id", telegramID))

	if res := s.str.Table("users").Create(&storage.User{
		TelegramID:  telegramID,
		LastSeen:    uint64(time.Now().Unix()),
		Coins:       coins,
		EarnedCoins: coins,
		Gold:        gold,
		Investors:   investors,
	}); res.Error != nil {
		return nil, res.Error
	}

	if user, err = s.SelectUser(telegramID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) SelectUser(telegramID uint64) (user *storage.User, err error) {
	s.lgr.Debug("selecting user", zap.Uint64("telegram_id", telegramID))

	res := s.str.Table("users").Where("telegram_id = ?", telegramID).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (s *Storage) SelectUsers() (users []storage.User, err error) {
	s.lgr.Debug("selecting all users")

	if res := s.str.Table("users").Find(&users); res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (s *Storage) UpdateUserCoins(telegramID, coins uint64) (user *storage.User, err error) {
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

func (s *Storage) UpdateUserGold(telegramID, gold uint64) (user *storage.User, err error) {
	s.lgr.Debug("updating user gold",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("gold", gold),
	)

	if res := s.str.Table("users").Where("telegram_id = ?", telegramID).Update("gold", gold); res.Error != nil {
		return nil, res.Error
	}

	if user, err = s.SelectUser(telegramID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) UpdateUserInvestors(telegramID, investors uint64) (user *storage.User, err error) {
	s.lgr.Debug("updating user investors",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("investors", investors),
	)

	if res := s.str.Table("users").Where("telegram_id = ?", telegramID).Update("investors", investors); res.Error != nil {
		return nil, res.Error
	}

	if user, err = s.SelectUser(telegramID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) UpdateUserEarnedCoins(telegramID, earnedCoins uint64) (user *storage.User, err error) {
	s.lgr.Debug("updating user earned coins",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("earned_coins", earnedCoins),
	)

	if res := s.str.Table("users").Where("telegram_id = ?", telegramID).Update("earned_coins", earnedCoins); res.Error != nil {
		return nil, res.Error
	}

	if user, err = s.SelectUser(telegramID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) UpdateUserLastSeen(telegramID, lastSeen uint64) (user *storage.User, err error) {
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

func (s *Storage) InsertUserCard(telegramID, cardID, level uint64) (userCard *storage.UserCard, err error) {
	s.lgr.Debug("inserting user card",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("card_id", cardID),
		zap.Uint64("level", level),
	)

	if res := s.str.Table("user_cards").Create(&storage.UserCard{
		TelegramID: telegramID,
		CardID:     cardID,
		Level:      level,
	}); res.Error != nil {
		return nil, res.Error
	}

	if userCard, err = s.SelectUserCard(telegramID, cardID); err != nil {
		return nil, err
	}

	return userCard, nil
}

func (s *Storage) SelectUserCard(telegramID, cardID uint64) (userCard *storage.UserCard, err error) {
	s.lgr.Debug("selecting user card",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("card_id", cardID),
	)

	if res := s.str.Table("user_cards").Where("telegram_id = ? AND card_id = ?", telegramID, cardID).First(&userCard); res.Error != nil {
		return nil, res.Error
	}

	return userCard, nil
}

func (s *Storage) SelectUserCards(telegramID uint64) (userCards []storage.UserCard, err error) {
	s.lgr.Debug("selecting user cards", zap.Uint64("telegram_id", telegramID))

	if res := s.str.Table("user_cards").Where("telegram_id = ?", telegramID).Find(&userCards); res.Error != nil {
		return nil, res.Error
	}

	return userCards, nil
}

func (s *Storage) UpdateUserCardLevel(telegramID, cardID, level uint64) (userCard *storage.UserCard, err error) {
	s.lgr.Debug("updating user card level",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("card_id", cardID),
		zap.Uint64("level", level),
	)

	if res := s.str.Table("user_cards").Where("telegram_id = ? AND card_id = ?", telegramID, cardID).Update("level", level); res.Error != nil {
		return nil, res.Error
	}

	if userCard, err = s.SelectUserCard(telegramID, cardID); err != nil {
		return nil, err
	}

	return userCard, nil
}

func (s *Storage) UpdateUserCardNextClick(telegramID, cardID, nextClick uint64) (userCard *storage.UserCard, err error) {
	s.lgr.Debug("updating user card next click",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("card_id", cardID),
		zap.Uint64("next_click", nextClick),
	)

	if res := s.str.Table("user_cards").Where("telegram_id = ? AND card_id = ?", telegramID, cardID).Update("next_click", nextClick); res.Error != nil {
		return nil, res.Error
	}

	if userCard, err = s.SelectUserCard(telegramID, cardID); err != nil {
		return nil, err
	}

	return userCard, nil
}

func (s *Storage) UpdateUserCardLastClick(telegramID, cardID, lastClick uint64) (userCard *storage.UserCard, err error) {
	s.lgr.Debug("updating user card last click",
		zap.Uint64("telegram_id", telegramID),
		zap.Uint64("card_id", cardID),
		zap.Uint64("last_click", lastClick),
	)

	if res := s.str.Table("user_cards").Where("telegram_id = ? AND card_id = ?", telegramID, cardID).Update("last_click", lastClick); res.Error != nil {
		return nil, res.Error
	}

	if userCard, err = s.SelectUserCard(telegramID, cardID); err != nil {
		return nil, err
	}

	return userCard, nil
}

func (s *Storage) SelectCard(cardID uint64) (cards *storage.Card, err error) {
	s.lgr.Debug("selecting card", zap.Uint64("card_id", cardID))

	if res := s.str.Table("cards").Where("id = ?", cardID).First(&cards); res.Error != nil {
		return nil, res.Error
	}

	return cards, nil
}

func (s *Storage) SelectCards() (cards []storage.Card, err error) {
	s.lgr.Debug("selecting all cards")

	if res := s.str.Table("cards").Find(&cards); res.Error != nil {
		return nil, res.Error
	}

	return cards, nil
}
