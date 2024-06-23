package model

type (
	User struct {
		ID         int    `json:"id"`
		TelegramID uint64 `json:"telegram_id"`
		Coins      uint64 `json:"coins"`
		LastSeen   uint64 `json:"last_seen"`
	}
)
