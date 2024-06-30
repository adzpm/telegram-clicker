package model

type (
	User struct {
		ID          uint64 `json:"id"`
		TelegramID  uint64 `json:"telegram_id"`
		LastSeen    uint64 `json:"last_seen"`
		Coins       uint64 `json:"coins"`
		EarnedCoins uint64 `json:"earned_coins"`
		Gold        uint64 `json:"gold"`
		Investors   uint64 `json:"investors"`
	}

	UserCard struct {
		ID         uint64 `json:"id"`
		TelegramID uint64 `json:"telegram_id"`
		CardID     uint64 `json:"card_id"`
		Level      uint64 `json:"level"`
		NextClick  uint64 `json:"next_click"`
		LastClick  uint64 `json:"last_click"`
	}

	Card struct {
		ID              uint64  `json:"id"`
		Name            string  `json:"name"`
		ImageURL        string  `json:"image_url"`
		Price           uint64  `json:"price"`
		PriceMultiplier float64 `json:"price_multiplier"`
		CoinsPerClick   uint64  `json:"coins_per_click"`
		ClickTimeout    uint64  `json:"click_timeout"`
		UpgradeLevel    uint64  `json:"upgrade_level"`
		MaxLevel        uint64  `json:"max_level"`
	}
)
