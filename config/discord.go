package config

import "fmt"

const (
	baseUrl	 = "https://discord.com/api"
	oauthUrl     = baseUrl + "/oauth2/authorize"
	redirect_uri = "http://localhost:8083/callback"
)

type DiscordConfig struct {
	ClientID     string `env:"DISCORD_CLIENT_ID" envDefault:""`
	ClientSecret string `env:"DISCORD_CLIENT_SECRET" envDefault:""`
}

func (d DiscordConfig) GetBaseURL() string {
	return baseUrl
}

func (d DiscordConfig) GetOAuthURL() string {
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=identify+email+guilds+guilds.join", oauthUrl, d.ClientID, redirect_uri)
}

func (d DiscordConfig) GetRedirectURI() string {
	return redirect_uri
}
