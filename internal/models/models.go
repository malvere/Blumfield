package models

type User struct {
	ID              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	LanguageCode    string `json:"language_code"`
	AllowsWriteToPM bool   `json:"allows_write_to_pm"`
}
type QueryData struct {
	QueryID  string `json:"query_id"`
	User     User   `json:"user"`
	AuthDate int64  `json:"auth_date"`
	Hash     string `json:"hash"`
}

type Query struct {
	Query string `json:"query"`
}

type BlumTokens struct {
	Token struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
		User    struct {
			ID struct {
				ID string `json:"id"`
			} `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"token"`
	JustCreated bool `json:"justCreated"`
}

type UnknownResponse struct {
	TotalFiatValue struct {
		Usd string `json:"usd"`
	} `json:"totalFiatValue"`
	Points []struct {
		CurrencyID string `json:"currencyId"`
		Name       string `json:"name"`
		Symbol     string `json:"symbol"`
		ImageURL   string `json:"imageUrl"`
		Balance    string `json:"balance"`
		Farming    struct {
			StartTime    int64  `json:"startTime"`
			EndTime      int64  `json:"endTime"`
			EarningsRate string `json:"earningsRate"`
			Balance      string `json:"balance"`
		} `json:"farming"`
		FiatValue struct {
		} `json:"fiatValue"`
	} `json:"points"`
}

type BalanceResponse struct {
	AvailableBalance     string `json:"availableBalance"`
	PlayPasses           int    `json:"playPasses"`
	IsFastFarmingEnabled bool   `json:"isFastFarmingEnabled"`
	Timestamp            int64  `json:"timestamp"`
	Farming              struct {
		StartTime    int64  `json:"startTime"`
		EndTime      int64  `json:"endTime"`
		EarningsRate string `json:"earningsRate"`
		Balance      string `json:"balance"`
	} `json:"farming"`
}

type MeResponse struct {
}
