package module

import (
	"context"

	"github.com/mocha-bot/mochus/core/entity"
	repository "github.com/mocha-bot/mochus/core/repository/discord"
)

type DiscordUsecase interface {
	ExchangeCodeForToken(ctx context.Context, code string) (string, error)
	GetUser(ctx context.Context, token string) (entity.User, error)
}

type discordUsecase struct {
	DiscordRepository repository.DiscordRepository
}

func NewDiscordUsecase(repo repository.DiscordRepository) DiscordUsecase {
	return &discordUsecase{
		DiscordRepository: repo,
	}
}

func (d *discordUsecase) ExchangeCodeForToken(ctx context.Context, code string) (string, error) {
	return d.DiscordRepository.GetTemporaryToken(ctx, code)
}

func (d *discordUsecase) GetUser(ctx context.Context, token string) (entity.User, error) {
	return d.DiscordRepository.GetUser(ctx, token)
}
