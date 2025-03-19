package discord_repository

import (
	"context"
	"net/http"

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

func (d *discordRepository) GetToken(ctx context.Context, code, requestURL string) (*entity.AccessToken, error) {
	var response AccessTokenResponse

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	payload := map[string]string{
		"grant_type":   GrantTypeAuthorizationCode,
		"code":         code,
		"redirect_uri": requestURL,
	}

	req := d.client.R().
		SetHeaders(headers).
		SetBasicAuth(d.cfg.ClientID, d.cfg.ClientSecret).
		SetContext(ctx).
		SetSuccessResult(&response).
		SetErrorResult(&response).
		SetFormData(payload)

	resp, err := req.Post(Oauth2GetToken)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, entity.ErrorUnauthorized
		}

		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("%w: %v", entity.ErrorBadRequest, response.HTTPResponse.Error)
		}

		return nil, fmt.Errorf("Error %v %v: %v", response.HTTPResponse.Message, response.HTTPResponse.Error, response.ErrorDescription)
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
	}

	req := d.client.R().
		SetHeaders(headers).
		SetBasicAuth(d.cfg.ClientID, d.cfg.ClientSecret).
		SetContext(ctx).
		SetSuccessResult(&response).
		SetErrorResult(&response).
		SetFormData(payload)

	resp, err := req.Post(Oauth2GetToken)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, entity.ErrorUnauthorized
		}

		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("%w: %v", entity.ErrorBadRequest, response.HTTPResponse.Error)
		}

		return nil, fmt.Errorf("Error %v %v: %v", response.HTTPResponse.Message, response.HTTPResponse.Error, response.ErrorDescription)
	}

	return response.AccessToken.ToEntity(), nil
}

func (d *discordRepository) RevokeToken(ctx context.Context, request *entity.RevokeTokenRequest) error {
	var response RevokeTokenResponse

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	req := d.client.R().
		SetHeaders(headers).
		SetBasicAuth(d.cfg.ClientID, d.cfg.ClientSecret).
		SetContext(ctx).
		SetFormData(request.ToPayload()).
		SetErrorResult(&response)

	resp, err := req.Post(Oauth2RevokeToken)
	if err != nil {
		return err
	}

	if resp.IsErrorState() {
		if resp.StatusCode == http.StatusUnauthorized {
			return entity.ErrorUnauthorized
		}

		if resp.StatusCode == http.StatusBadRequest {
			return fmt.Errorf("%w: %v", entity.ErrorBadRequest, response.HTTPResponse.Error)
		}

		return fmt.Errorf("Error %v %v: %v", response.HTTPResponse.Message, response.HTTPResponse.Error, response.HTTPResponse.ErrorDescription)
	}

	return nil
}

func (d *discordRepository) GetUser(ctx context.Context, token string) (*entity.User, error) {
	var response UserResponse

	req := d.client.R().
		SetHeader("Authorization", token).
		SetContext(ctx).
		SetSuccessResult(&response).
		SetErrorResult(&response)

	resp, err := req.Get(GetUser)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, entity.ErrorUnauthorized
		}

		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("%w: %v", entity.ErrorBadRequest, response.HTTPResponse.Error)
		}

		return nil, fmt.Errorf("Error %v %v: %v", response.HTTPResponse.Message, response.HTTPResponse.Error, response.ErrorDescription)
	}

	return response.User.ToEntity(), nil
}
