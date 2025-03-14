package module

import (
	"context"
	"fmt"

	"github.com/mocha-bot/mochus/core/entity"
	repository "github.com/mocha-bot/mochus/core/repository/discord"
)

type DiscordUsecase interface {
	ExchangeCodeForToken(ctx context.Context, code string) (*entity.AccessToken, error)
	ExchangeRefreshForToken(ctx context.Context, refreshToken string) (*entity.AccessToken, error)
	RevokeToken(ctx context.Context, token string) error

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

func (d *discordUsecase) ExchangeRefreshForToken(ctx context.Context, refreshToken string) (*entity.AccessToken, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("%w, invalid refresh token", entity.ErrorUnauthorized)
	}

	accessToken, err := d.DiscordRepository.GetTokenByRefresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
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

func (d *discordUsecase) RevokeToken(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("%w, invalid token", entity.ErrorUnauthorized)
	}

	return d.DiscordRepository.RevokeToken(ctx, token)
}
