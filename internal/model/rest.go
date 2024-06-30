package model

type (
	Game struct {
		UserID                        uint64               `json:"user_id"`
		TelegramID                    uint64               `json:"telegram_id"`
		LastSeen                      uint64               `json:"last_seen"`
		CurrentCoins                  uint64               `json:"current_coins"`
		CurrentGold                   uint64               `json:"current_gold"`
		CurrentInvestors              uint64               `json:"current_investors"`
		InvestorsAfterReset           uint64               `json:"investors_after_reset"`
		CurrentInvestorsMultiplier    float64              `json:"current_investors_multiplier"`
		InvestorsMultiplierAfterReset float64              `json:"investors_multiplier_after_reset"`
		PercentsPerInvestor           uint64               `json:"percents_per_investor"`
		Cards                         map[uint64]*GameCard `json:"cards"`
	}

	GameCard struct {
		ID       uint64 `json:"id"`
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`

		CurrentLevel uint64 `json:"current_level"`
		MaxLevel     uint64 `json:"max_level"`

		CurrentPrice   uint64 `json:"current_price"`
		NextLevelPrice uint64 `json:"upgrade_price"`

		ClickTimeout uint64 `json:"click_timeout"`
		NextClick    uint64 `json:"next_click"`
		LastClick    uint64 `json:"last_click"`

		CurrentCoinsPerClick   uint64 `json:"current_coins_per_click"`
		NextLevelCoinsPerClick uint64 `json:"next_level_coins_per_click"`
	}
)
