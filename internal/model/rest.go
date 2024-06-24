package model

type (
	Game struct {
		UserID       uint64                  `json:"user_id"`
		TelegramID   uint64                  `json:"telegram_id"`
		LastSeen     uint64                  `json:"last_seen"`
		CurrentCoins uint64                  `json:"current_coins"`
		CurrentGold  uint64                  `json:"current_gold"`
		Products     map[uint64]*GameProduct `json:"products"`
	}

	GameProduct struct {
		ID             uint64 `json:"id"`
		Name           string `json:"name"`
		ImageURL       string `json:"image_url"`
		UpgradePrice   uint64 `json:"upgrade_price"`
		CoinsPerClick  uint64 `json:"coins_per_click"`
		CoinsPerMinute uint64 `json:"coins_per_minute"`
		CurrentLevel   uint64 `json:"current_level"`
		MaxLevel       uint64 `json:"max_level"`
	}
)
