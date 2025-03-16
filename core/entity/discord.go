package entity

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	cookiey "github.com/mocha-bot/mochus/pkg/cookie"
)

const (
	OAuthRefreshTokenMaxAge = time.Hour
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

	refreshTokenMaxAge := at.ExpiresIn + int(OAuthRefreshTokenMaxAge.Seconds())

	accessTokenCookie := &http.Cookie{
		Name:   cookiey.CookieAccessToken,
		Value:  at.AccessToken,
		MaxAge: at.ExpiresIn,
	}

	refreshTokenCookie := &http.Cookie{
		Name:   cookiey.CookieRefreshToken,
		Value:  at.RefreshToken,
		MaxAge: refreshTokenMaxAge,
	}

	tokenTypeCookie := &http.Cookie{
		Name:   cookiey.CookieTokenType,
		Value:  at.TokenType,
		MaxAge: at.ExpiresIn,
	}

	scopeCookie := &http.Cookie{
		Name:   cookiey.CookieScope,
		Value:  at.Scope,
		MaxAge: at.ExpiresIn,
	}

	return []*http.Cookie{accessTokenCookie, refreshTokenCookie, tokenTypeCookie, scopeCookie}
}

type RefreshTokenRequest struct {
	RefreshToken string `validate:"required"`
	Referer      string `header:"Referer" validate:"required"`
}

func (rt *RefreshTokenRequest) Validate() error {
	err := validator.New().Struct(rt)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field %s is invalid", e.Field())
		}
	}

	return nil
}

const (
	GrantTypeRefreshToken = "refresh_token"
	GrantTypeAccessToken  = "access_token"
)

type RevokeTokenRequest struct {
	AccessToken  string // cookie
	RefreshToken string // cookie
	Referer      string `header:"Referer" validate:"required"`
}

func (rt *RevokeTokenRequest) Validate() error {
	err := validator.New().Struct(rt)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field %s is invalid", e.Field())
		}
	}

	return nil
}

func (rt *RevokeTokenRequest) ToPayload() map[string]string {
	payload := make(map[string]string)

	switch {
	case rt.RefreshToken != "":
		payload["token"] = rt.RefreshToken
		payload["token_type_hint"] = GrantTypeRefreshToken
	case rt.AccessToken != "":
		payload["token"] = rt.AccessToken
		payload["token_type_hint"] = GrantTypeAccessToken
	}

	return payload
}

type GetUserByTokenRequest struct {
	Authorization string `header:"Authorization"`
}
