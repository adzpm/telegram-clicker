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

	UserProduct struct {
		ID         uint64 `json:"id"`
		TelegramID uint64 `json:"telegram_id"`
		ProductID  uint64 `json:"product_id"`
		Level      uint64 `json:"level"`
	}

	Product struct {
		ID                      uint64  `json:"id"`
		Name                    string  `json:"name"`
		ImageURL                string  `json:"image_url"`
		StartProductPrice       uint64  `json:"start_product_price"`
		ProductPriceMultiplier  float64 `json:"product_price_multiplier"`
		StartCoinsPerClick      uint64  `json:"start_coins_per_click"`
		CoinsPerClickMultiplier float64 `json:"coins_per_click_multiplier"`
		MaxLevel                uint64  `json:"max_level"`
	}
)
