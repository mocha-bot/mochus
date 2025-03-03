package repository

import (
	"context"

	"github.com/mocha-bot/mochus/core/entity"
)

type DiscordRepository interface {
	GetTemporaryToken(ctx context.Context, code string) (string, error)
	GetUser(ctx context.Context, token string) (entity.User, error)
}
