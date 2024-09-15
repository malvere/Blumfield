package blumfield

import (
	"blumfield/internal/models"
	"context"
	"math/rand/v2"
	"time"
)

//PLAY : /play
//CLAIM : /claim

func (b *Blumfield) PlayGame(ctx context.Context) {
	balance, err := b.GetBalance()
	if err != nil {
		return
	}
	b.log.Info("Available play-passes: ", balance.PlayPasses)
	if balance.PlayPasses <= 1 {
		b.log.Info("Insufficient ticket balance.")
	}
	for {
		select {
		case <-ctx.Done():
			b.log.Info("Shutting down...")
			return
		default:
			balance, err := b.GetBalance()
			if err != nil {
				b.log.Error("Error retrieving balance...")
				continue
			}
			b.log.Info("Available balance: ", balance.AvailableBalance)
			b.log.Info("Available play-passes: ", balance.PlayPasses)
			if balance.PlayPasses <= 1 {
				b.log.Info("Insufficient ticket balance.")
				return
			}
		}
		var game *models.PlayGameResponse

		b.log.Info("Starting game...")
		start, err := b.client.R().
			SetHeaders(b.BaseHeaders).
			SetResult(&models.PlayGameResponse{}).
			Post(gameURL + "/play")
		if err != nil {
			b.log.Error("Error obtaining gameID: ", err)
			b.log.Error(start.String())
		}
		game = start.Result().(*models.PlayGameResponse)

		b.log.Info("Waiting 30s...")
		time.Sleep(time.Duration(rand.IntN(10)+32) * time.Second)
		points := rand.IntN(50) + 200

		b.log.Info("Claiming reward...")
		claim, err := b.client.R().
			SetHeaders(b.BaseHeaders).
			SetBody(&models.ClaimGameRequest{
				GameID: game.GameID,
				Points: points,
			}).
			Post(gameURL + "/claim")
		if err != nil {
			b.log.Error("Error claiming game reward: ", err)
			b.log.Error(claim.String())
		}
		if claim.String() == "OK" {
			b.log.Info("Succesfully claimed ", points, " points!")
		} else {
			return
		}
	}
}
