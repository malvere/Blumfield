package models

type PlayGameResponse struct {
	GameID string `json:"gameId"`
}

type ClaimGameRequest struct {
	GameID string `json:"gameId"`
	Points int    `json:"points"`
}

// Claim response "OK"
