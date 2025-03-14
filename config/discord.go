package config

import "fmt"

const (
	DiscordBaseURL  = "https://discord.com/api"
	DiscordOAuthURL = DiscordBaseURL + "/oauth2/authorize"
)

type DiscordConfig struct {
	ClientID      string `env:"DISCORD_CLIENT_ID" envDefault:""`
	ClientSecret  string `env:"DISCORD_CLIENT_SECRET" envDefault:""`
	RedirectURI   string `env:"DISCORD_REDIRECT_URI" envDefault:""`
	LatestVersion string `env:"DISCORD_LATEST_VERSION" envDefault:"v10"`
}

func (d DiscordConfig) GetBaseURL(version string) string {
	if version == "" {
		version = d.LatestVersion
	}

	return fmt.Sprintf("%s/%s", DiscordBaseURL, version)
}

func (d DiscordConfig) GetRedirectURI() string {
	return d.RedirectURI
}
