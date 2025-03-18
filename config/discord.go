package config

import "fmt"

type DiscordConfig struct {
	ClientID       string `env:"DISCORD_CLIENT_ID" envDefault:""`
	ClientSecret   string `env:"DISCORD_CLIENT_SECRET" envDefault:""`
	BaseURL        string `env:"DISCORD_BASE_URL" envDefault:"https://discord.com/api"`
	RedirectDomain string `env:"DISCORD_REDIRECT_DOMAIN" envDefault:".mocha-bot.xyz"`
	LatestVersion  string `env:"DISCORD_LATEST_VERSION" envDefault:"v10"`
}

func (d DiscordConfig) GetBaseURL(version string) string {
	if version == "" {
		version = d.LatestVersion
	}

	return fmt.Sprintf("%s/%s", d.BaseURL, version)
}
