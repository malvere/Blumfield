package blumfield

import (
	"blumfield/internal/models"
	"context"
)

func (b *Blumfield) ClaimFarming() (*models.ClaimFarmingResponse, error) {
	resp, err := b.client.R().SetHeaders(b.BaseHeaders).SetResult(&models.ClaimFarmingResponse{}).Post(farmingURL + "/claim")
	if err != nil {
		return nil, err
	}
	return resp.Result().(*models.ClaimFarmingResponse), nil
}

func (b *Blumfield) StartFarming() (*models.StartFarmingResponse, error) {
	resp, err := b.client.R().
		SetHeaders(b.BaseHeaders).
		SetResult(&models.StartFarmingResponse{}).
		Post(farmingURL + "/start")
	if err != nil {
		return nil, err
	}
	return resp.Result().(*models.StartFarmingResponse), nil

}

func (b *Blumfield) ClaimCheckIn() error {
	respChek, err := b.client.R().
		SetHeaders(b.BaseHeaders).
		Get(checkInURL)
	if err != nil {
		return err
	}

	if respChek.StatusCode() == 404 {
		b.log.Info("Alredy checked in today")
		return nil
	}

	respClaim, err := b.client.R().SetHeaders(b.BaseHeaders).Post(checkInURL)
	if err != nil {
		return err
	}

	if respClaim.String() != "OK" {
		b.log.Info("Check-in already complete", respClaim.String())
		return nil
	}
	b.log.Info("Checked in succesfully!")
	return nil
}

func (b *Blumfield) farmWithContext(ctx context.Context) {
	claimed := false // Track whether farming was claimed

	// Check context before starting
	if ctx.Err() != nil {
		b.log.Info("Context canceled, stopping farmWithContext")
		return
	}

	if err := b.LogBalance(ctx); err != nil {
		b.log.Error("Error checking balance: ", err)
	}

	// Delay with context check
	if !b.tools.DelayWithContext(ctx, 3) {
		b.log.Info("Context canceled during delay, stopping farmWithContext")
		return
	}

	// Claim farming
	farmClaim, err := b.ClaimFarming()
	if err != nil {
		b.log.Error("Error claiming farming: ", err)
	} else {
		claimed = true // Mark that farming was claimed
		if farmClaim.AvailableBalance == "" {
			b.log.Info("Already claimed.")
		} else {
			b.log.Info("Claimed! Balance: ", farmClaim.AvailableBalance)
		}
	}

	// Check context before proceeding, but ensure StartFarming is executed if claimed
	if ctx.Err() != nil && !claimed {
		b.log.Info("Context canceled before starting farming")
		return
	}

	// Log balance before starting farming
	if err := b.LogBalance(ctx); err != nil {
		b.log.Error("Error checking balance: ", err)
	}

	// Delay with context check, but ensure StartFarming is executed if claimed
	if !b.tools.DelayWithContext(ctx, 2) && !claimed {
		b.log.Info("Context canceled during delay, stopping farmWithContext")
		return
	}

	// Start farming (must execute if ClaimFarming was executed)
	farmStart, err := b.StartFarming()
	if err != nil {
		b.log.Error("Error starting farming: ", err)
	} else {
		b.log.Info("Farming: ", farmStart.EarningsRate)
	}
}
