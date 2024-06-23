package model

type (
	GameInformation struct {
		ID           uint64         `json:"id"`
		UserID       uint64         `json:"user_id"`
		LastSeen     uint64         `json:"last_seen"`
		Coins        uint64         `json:"coins"`
		Products     []*Product     `json:"products"`
		UserProducts []*UserProduct `json:"user_products"`
	}
)
