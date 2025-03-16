package module

import (
	"context"
	"fmt"

	"github.com/mocha-bot/mochus/core/entity"
	repository "github.com/mocha-bot/mochus/core/repository/discord"
)

type DiscordUsecase interface {
	ExchangeCodeForToken(ctx context.Context, code string) (*entity.AccessToken, error)
	ExchangeRefreshForToken(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.AccessToken, error)
	RevokeToken(ctx context.Context, req *entity.RevokeTokenRequest) error

	GetUser(ctx context.Context, token string) (*entity.User, error)
}

type discordUsecase struct {
	DiscordRepository repository.DiscordRepository
}

func NewDiscordUsecase(repo repository.DiscordRepository) DiscordUsecase {
	return &discordUsecase{
		DiscordRepository: repo,
	}
}

func (d *discordUsecase) ExchangeCodeForToken(ctx context.Context, code string) (*entity.AccessToken, error) {
	if code == "" {
		return nil, fmt.Errorf("%w, invalid code", entity.ErrorUnauthorized)
	}

	accessToken, err := d.DiscordRepository.GetToken(ctx, code)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (d *discordUsecase) ExchangeRefreshForToken(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.AccessToken, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w, %w", entity.ErrorUnauthorized, err)
	}

	accessToken, err := d.DiscordRepository.GetTokenByRefresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (d *discordUsecase) RevokeToken(ctx context.Context, req *entity.RevokeTokenRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w, %w", entity.ErrorUnauthorized, err)
	}

	return d.DiscordRepository.RevokeToken(ctx, req)
}

func (d *discordUsecase) GetUser(ctx context.Context, token string) (*entity.User, error) {
	if token == "" {
		return nil, fmt.Errorf("%w, invalid token", entity.ErrorUnauthorized)
	}

	user, err := d.DiscordRepository.GetUser(ctx, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}
