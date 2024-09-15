package blumfield

import "blumfield/internal/models"

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

func (b *Blumfield) LogBalance() error {
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
