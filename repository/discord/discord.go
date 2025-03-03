package discord_repository

import (
	"context"

	// "encoding/base64"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/mocha-bot/mochus/config"
	"github.com/mocha-bot/mochus/core/entity"
	repository "github.com/mocha-bot/mochus/core/repository/discord"
)

type discordRepository struct {
	client *req.Client
	cfg    config.DiscordConfig
}

func NewDiscordRepository(cfg config.DiscordConfig) repository.DiscordRepository {
	return &discordRepository{
		client: req.NewClient().SetBaseURL(cfg.GetBaseURL()),
		cfg:    cfg,
	}
}

func (d *discordRepository) GetTemporaryToken(ctx context.Context, code string) (string, error) {
	var response TokenResponse
	var errorResponse ErrorResponse
	req := d.client.R()
	resp, err := req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}).
		SetBasicAuth(d.cfg.ClientID, d.cfg.ClientSecret).
		SetContext(ctx).
		SetSuccessResult(&response).
		SetErrorResult(&errorResponse).
		SetFormData(map[string]string{
			"grant_type":   "authorization_code",
			"code":         code,
			"redirect_uri": d.cfg.GetRedirectURI(),
		}).
		Post("/v10/oauth2/token")

	if err != nil {
		return "", err
	}

	if resp.IsErrorState() {
		return "", fmt.Errorf("Error with status code: %v", resp.StatusCode)
	}

	return response.AccessToken, nil
}

func (d *discordRepository) GetUser(ctx context.Context, token string) (entity.User, error) {
	var response entity.User

	req := d.client.R()
	resp, err := req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetContext(ctx).
		SetSuccessResult(&response).
		Get("/users/@me")

	if err != nil {
		return entity.User{}, err
	}

	if resp.IsErrorState() {
		return entity.User{}, fmt.Errorf("Error with status code: %v", resp.StatusCode)
	}

	return response, nil
}
