package model

type (
	Game struct {
		UserID       uint64                  `json:"user_id"`
		TelegramID   uint64                  `json:"telegram_id"`
		LastSeen     uint64                  `json:"last_seen"`
		CurrentCoins uint64                  `json:"current_coins"`
		Products     map[uint64]*GameProduct `json:"products"`
	}

	GameProduct struct {
		ID       uint64 `json:"id"`
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`

		CurrentLevel uint64 `json:"current_level"`
		MaxLevel     uint64 `json:"max_level"`

		CurrentPrice   uint64 `json:"current_price"`
		NextLevelPrice uint64 `json:"upgrade_price"`

		CurrentCoinsPerClick   uint64 `json:"current_coins_per_click"`
		NextLevelCoinsPerClick uint64 `json:"next_level_coins_per_click"`
	}
)
