package entity

import (
	"net/http"
	"strings"
	"time"
)

type OauthCallbackRequest struct {
	Code string `query:"code"`
}

type AccessToken struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	RefreshToken string
	Scope        string
}

func (at *AccessToken) ToHTTPCookies() Cookies {
	if at == nil {
		return nil
	}

	standardExpires := time.Now().Add(time.Duration(at.ExpiresIn) * time.Second)
	refreshTokenExpires := standardExpires.Add(1 * time.Hour)

	accessTokenCookie := &http.Cookie{
		Name:   "access_token",
		Value:  at.AccessToken,
		MaxAge: at.ExpiresIn,
	}

	refreshTokenCookie := &http.Cookie{
		Name:   "refresh_token",
		Value:  at.RefreshToken,
		MaxAge: refreshTokenExpires.Second(),
	}

	tokenTypeCookie := &http.Cookie{
		Name:   "token_type",
		Value:  at.TokenType,
		MaxAge: at.ExpiresIn,
	}

	scopeCookie := &http.Cookie{
		Name:   "scope",
		Value:  at.Scope,
		MaxAge: at.ExpiresIn,
	}

	return []*http.Cookie{accessTokenCookie, refreshTokenCookie, tokenTypeCookie, scopeCookie}
}

type RefreshTokenRequest struct {
	RefreshToken string `header:"refresh_token"`
}

type RevokeTokenRequest struct {
	Token string `header:"token"`
}

type GetUserByTokenRequest struct {
	Authorization string `header:"Authorization"`
}

func (r *GetUserByTokenRequest) Token() string {
	strs := strings.Split(r.Authorization, " ")
	if len(strs) != 2 {
		return ""
	}

	if strs[0] != "Bearer" {
		return ""
	}

	return strs[1]
}
