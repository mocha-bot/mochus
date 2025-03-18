package discord_repository

import (
	"github.com/mocha-bot/mochus/core/entity"
	"github.com/mocha-bot/mochus/pkg/discord"
)

const (
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeClientCredentials = "client_credentials"
	GrantTypeRefreshToken      = "refresh_token"
	GrantTypeAccessToken       = "access_token"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func (t *AccessToken) ToEntity() *entity.AccessToken {
	if t == nil {
		return nil
	}

	return &entity.AccessToken{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		ExpiresIn:    t.ExpiresIn,
		Scope:        t.Scope,
		RefreshToken: t.RefreshToken,
	}
}

type AccessTokenPayload struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

type AccessTokenResponse struct {
	*AccessToken
	*discord.HTTPResponse
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	PublicFlags   int    `json:"public_flags"`
	Flags         int    `json:"flags"`
	Banner        string `json:"banner"`
	AccentColor   string `json:"accent_color,omitempty"`
	GlobalName    string `json:"global_name"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Locale        string `json:"locale"`
	PremiumType   int    `json:"premium_type"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
}

func (u *User) ToEntity() *entity.User {
	if u == nil {
		return nil
	}

	return &entity.User{
		ID:            u.ID,
		Username:      u.Username,
		Discriminator: u.Discriminator,
		Avatar:        u.Avatar,
		PublicFlags:   u.PublicFlags,
		Flags:         u.Flags,
		Banner:        u.Banner,
		AccentColor:   u.AccentColor,
		GlobalName:    u.GlobalName,
		MFAEnabled:    u.MFAEnabled,
		Locale:        u.Locale,
		PremiumType:   u.PremiumType,
		Email:         u.Email,
		Verified:      u.Verified,
	}
}

type UserResponse struct {
	*User
	*discord.HTTPResponse
}

type RevokeTokenResponse struct {
	*discord.HTTPResponse
}
