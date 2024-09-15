package models

type DailyGetResponse struct {
	Days []Day `json:"days"`
}

type Day struct {
	Ordinal int `json:"ordinal"`
	Reward  struct {
		Passes int    `json:"passes"`
		Points string `json:"points"`
	} `json:"reward"`
}

// POST Response is "OK"
