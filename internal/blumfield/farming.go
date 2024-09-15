package blumfield

import "blumfield/internal/models"

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
