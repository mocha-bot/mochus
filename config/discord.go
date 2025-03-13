package config

import "fmt"

const (
	DiscordBaseURL  = "https://discord.com/api"
	DiscordOAuthURL = DiscordBaseURL + "/oauth2/authorize"
)

type DiscordConfig struct {
	ClientID     string `env:"DISCORD_CLIENT_ID" envDefault:""`
	ClientSecret string `env:"DISCORD_CLIENT_SECRET" envDefault:""`
	RedirectURI  string `env:"DISCORD_REDIRECT_URI" envDefault:""`
}

func (d DiscordConfig) GetBaseURL() string {
	return DiscordBaseURL
}

func (d DiscordConfig) GetOAuthURL() string {
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=identify+email+guilds+guilds.join", DiscordOAuthURL, d.ClientID, d.GetRedirectURI())
}

func (d DiscordConfig) GetRedirectURI() string {
	return d.RedirectURI
}
