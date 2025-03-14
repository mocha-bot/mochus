package discord_repository

import (
	"context"

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
		client: req.NewClient().SetBaseURL(cfg.GetBaseURL(cfg.LatestVersion)),
		cfg:    cfg,
	}
}

func (d *discordRepository) GetToken(ctx context.Context, code string) (*entity.AccessToken, error) {
	var response AccessTokenResponse

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	payload := map[string]string{
		"grant_type":   GrantTypeAuthorizationCode,
		"code":         code,
		"redirect_uri": d.cfg.GetRedirectURI(),
	}

	req := d.client.R().
		SetHeaders(headers).
		SetBasicAuth(d.cfg.ClientID, d.cfg.ClientSecret).
		SetContext(ctx).
		SetSuccessResult(&response).
		SetErrorResult(&response).
		SetFormData(payload)

	resp, err := req.Post("/oauth2/token")
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("Error %v: %v", response.HTTPResponse.Error, response.ErrorDescription)
	}

	return response.AccessToken.ToEntity(), nil
}

func (d *discordRepository) GetTokenByRefresh(ctx context.Context, refreshToken string) (*entity.AccessToken, error) {
	var response AccessTokenResponse

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	payload := map[string]string{
		"grant_type":    GrantTypeRefreshToken,
		"refresh_token": refreshToken,
		"redirect_uri":  d.cfg.GetRedirectURI(),
	}

	req := d.client.R().
		SetHeaders(headers).
		SetBasicAuth(d.cfg.ClientID, d.cfg.ClientSecret).
		SetContext(ctx).
		SetSuccessResult(&response).
		SetErrorResult(&response).
		SetFormData(payload)

	resp, err := req.Post("/oauth2/token")
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("Error %v: %v", response.HTTPResponse.Error, response.HTTPResponse.ErrorDescription)
	}

	return response.AccessToken.ToEntity(), nil
}

func (d *discordRepository) RevokeToken(ctx context.Context, token string) error {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	payload := map[string]string{
		"token":           token,
		"token_type_hint": GrantTypeAccessToken,
	}

	req := d.client.R().
		SetHeaders(headers).
		SetBasicAuth(d.cfg.ClientID, d.cfg.ClientSecret).
		SetContext(ctx).
		SetFormData(payload)

	resp, err := req.Post("/oauth2/token/revoke")
	if err != nil {
		return err
	}

	if resp.IsErrorState() {
		return fmt.Errorf("Error with status code: %v", resp.StatusCode)
	}

	return nil
}

func (d *discordRepository) GetUser(ctx context.Context, token string) (*entity.User, error) {
	var response UserResponse

	req := d.client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetContext(ctx).
		SetSuccessResult(&response)

	resp, err := req.Get("/users/@me")
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, fmt.Errorf("Error %v: %v", response.HTTPResponse.Error, response.ErrorDescription)
	}

	return response.User.ToEntity(), nil
}
