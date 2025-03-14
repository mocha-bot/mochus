package discord_repository

import (
	"github.com/mocha-bot/mochus/core/entity"
	"github.com/mocha-bot/mochus/pkg/discord"
)

const (
	GrantTypeAuthorizationCode = "authorization_code"
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
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
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
		Email:         u.Email,
	}
}

type UserResponse struct {
	*User
	*discord.HTTPResponse
}
