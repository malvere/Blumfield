package blumfield

import (
	jwts "blumfield/internal/jwt"
	"blumfield/internal/models"
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

func (b *Blumfield) PlayGame(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			b.log.Info("Shutting down...")
			return
		default:
			if err := b.ensureValidToken(); err != nil {
				b.log.Error("Token is invalid or expired, cannot renew: ", err)
				return
			}
			balance, err := b.GetBalance()
			if err != nil {
				b.log.Error("Error retrieving balance...")
				return
			}
			if balance.PlayPasses <= 1 {
				b.log.Info("Insufficient ticket balance.")
				return
			}
			b.log.Info("Available balance: ", balance.AvailableBalance)
			b.log.Info("Available play-passes: ", balance.PlayPasses)
			if err := b.startAndClaimGame(balance); err != nil {
				b.log.Error("Error during game play: ", err)
				return
			}
		}
	}
}

func (b *Blumfield) ensureValidToken() error {
	ttl, err := jwts.GetTokenEXP(b.Tokens.Token.Access)
	if err != nil || ttl <= 2*time.Minute {
		b.log.Warn("Token invalid or expiring soon, renewing...")
		if err := b.RenewAccessToken(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Blumfield) startAndClaimGame(balance *models.BalanceResponse) error {
	gameCtx := context.Background()
	b.tools.LogGameStatus(balance, "Starting...")
	var game *models.PlayGameResponse
	start, err := b.client.R().
		SetHeaders(b.BaseHeaders).
		SetContext(gameCtx).
		SetResult(&models.PlayGameResponse{}).
		Post(gameURL + "/play")
	if err != nil {
		b.log.Error("Error obtaining gameID: ", err)
		b.log.Error(start.String())
		return err
	}
	game = start.Result().(*models.PlayGameResponse)
	b.tools.LogGameStatus(balance, "Waiting 30s...")

	time.Sleep(time.Duration(rand.IntN(10)+32) * time.Second)
	points := rand.IntN(50) + 200
	claim, err := b.client.R().
		SetContext(gameCtx).
		SetHeaders(b.BaseHeaders).
		SetBody(&models.ClaimGameRequest{
			GameID: game.GameID,
			Points: points,
		}).
		Post(gameURL + "/claim")
	if err != nil || claim.String() != "OK" {
		b.log.Error("Error claiming game reward: ", err)
		b.log.Error(claim.String())
		return err
	}
	success := fmt.Sprint("Claimed ", points, " points!")
	b.tools.LogGameStatus(balance, success)
	return nil
}
