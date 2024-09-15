package blumfield

import "fmt"

const (
	tokensURL  = "https://user-domain.blum.codes/api/v1/auth/provider/PROVIDER_TELEGRAM_MINI_APP"
	balanceURL = "https://game-domain.blum.codes/api/v1/user/balance"
	earnURL    = "https://earn-domain.blum.codes/api/v1/tasks"
	farmingURL = "https://game-domain.blum.codes/api/v1/farming"
	checkInURL = "https://game-domain.blum.codes/api/v1/daily-reward?offset=-180"
	gameURL    = "https://game-domain.blum.codes/api/v1/game"
)

func StartTaskURL(taskID string) string {
	return fmt.Sprintf("%s/%s/start", earnURL, taskID)
}

func ClaimTaskURL(taskID string) string {
	return fmt.Sprintf("%s/%s/claim", earnURL, taskID)
}
