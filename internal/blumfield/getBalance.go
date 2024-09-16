package blumfield

import (
	"blumfield/internal/models"
	"context"
)

func (b *Blumfield) GetBalance() (*models.BalanceResponse, error) {
	resp, err := b.client.R().
		SetHeaders(b.BaseHeaders).
		SetResult(&models.BalanceResponse{}).
		Get(balanceURL)

	if err != nil {
		return nil, err
	}
	return resp.Result().(*models.BalanceResponse), nil
}

func (b *Blumfield) LogBalance(ctx context.Context) error {
	select {
	case <-ctx.Done():
		b.log.Info("Context canceled, stopping LogBalance")
		return ctx.Err()
	default:
		resp, err := b.client.R().
			SetHeaders(b.BaseHeaders).
			SetResult(&models.BalanceResponse{}).
			Get(balanceURL)

		if err != nil {
			return err
		}
		b.log.Info("Balance: ", resp.Result().(*models.BalanceResponse).AvailableBalance)
		return nil
	}
}
