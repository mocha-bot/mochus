package repository

import (
	"context"

	"github.com/mocha-bot/mochus/core/entity"
)

type DiscordRepository interface {
	GetToken(ctx context.Context, code string) (*entity.AccessToken, error)
	GetTokenByRefresh(ctx context.Context, refreshToken string) (*entity.AccessToken, error)
	RevokeToken(ctx context.Context, token string) error

	GetUser(ctx context.Context, token string) (*entity.User, error)
}
