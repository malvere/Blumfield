package models

type ClaimFarmingResponse struct {
	AvailableBalance     string `json:"availableBalance"`
	PlayPasses           int    `json:"playPasses"`
	IsFastFarmingEnabled bool   `json:"isFastFarmingEnabled"`
	Timestamp            int64  `json:"timestamp"`
}

type StartFarmingResponse struct {
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	EarningsRate string `json:"earningsRate"`
	Balance      string `json:"balance"`
}
