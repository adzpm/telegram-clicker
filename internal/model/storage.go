package model

type (
	User struct {
		ID         uint64 `json:"id"`
		TelegramID uint64 `json:"telegram_id"`
		LastSeen   uint64 `json:"last_seen"`
		Coins      uint64 `json:"coins"`
	}

	UserProduct struct {
		ID         uint64 `json:"id"`
		TelegramID uint64 `json:"telegram_id"`
		ProductID  uint64 `json:"product_id"`
		Level      uint64 `json:"level"`
	}

	Product struct {
		ID       uint64 `json:"id"`
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`

		StartPrice      uint64  `json:"start_price"`
		PriceMultiplier float64 `json:"price_multiplier"`

		StartCoins      uint64  `json:"start_coins"`
		CoinsMultiplier float64 `json:"coins_multiplier"`

		MaxLevel uint64 `json:"max_level"`
	}
)
