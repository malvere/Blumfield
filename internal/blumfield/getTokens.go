package blumfield

import (
	"blumfield/config"
	"blumfield/internal/models"
	"fmt"
)

func (b *Blumfield) RenewAccessToken() error {
	request := models.Query{
		Query: b.config.Auth.WebAppInit,
	}

	resp, err := b.client.R().
		SetBody(request).
		SetHeaders(b.BaseHeaders).
		SetResult(&models.BlumTokens{}).
		Post(tokensURL)

	if err != nil {
		fmt.Print("Error making request request: ", err)
		return err
	}

	// Extract token
	token := resp.Result().(*models.BlumTokens)
	b.BaseHeaders["Authorization"] = "Bearer " + token.Token.Access
	b.Tokens = token
	b.config.SaveTokens(&config.Tokens{
		Auth:    token.Token.Access,
		Refresh: token.Token.Refresh,
	})

	return nil
}
