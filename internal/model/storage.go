package model

type (
	User struct {
		ID         uint64 `json:"id"`
		TelegramID uint64 `json:"telegram_id"`
		LastSeen   uint64 `json:"last_seen"`
		Coins      uint64 `json:"coins"`
		Gold       uint64 `json:"gold"`
	}

	UserProduct struct {
		ID         uint64 `json:"id"`
		TelegramID uint64 `json:"telegram_id"`
		ProductID  uint64 `json:"product_id"`
		Level      uint64 `json:"level"`
	}

	Product struct {
		ID                       uint64  `json:"id"`
		Name                     string  `json:"name"`
		ImageURL                 string  `json:"image_url"`
		StartPrice               uint64  `json:"start_price"`
		PriceMultiplier          float64 `json:"price_multiplier"`
		StartCoinsPerClick       uint64  `json:"start_coins_per_click"`
		CoinsMultiplier          float64 `json:"coins_multiplier"`
		UpgradeProductLevel      uint64  `json:"upgrade_product_level"`
		UpgradeProductMultiplier float64 `json:"upgrade_product_multiplier"`
		MaxLevel                 uint64  `json:"max_level"`
	}
)
