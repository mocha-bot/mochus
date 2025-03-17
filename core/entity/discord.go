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
		Name:     cookiey.CookieAccessToken,
		Value:    at.AccessToken,
		MaxAge:   at.ExpiresIn,
		HttpOnly: true,
	}

	refreshTokenCookie := &http.Cookie{
		Name:     cookiey.CookieRefreshToken,
		Value:    at.RefreshToken,
		MaxAge:   refreshTokenMaxAge,
		HttpOnly: true,
	}

	tokenTypeCookie := &http.Cookie{
		Name:     cookiey.CookieTokenType,
		Value:    at.TokenType,
		MaxAge:   at.ExpiresIn,
		HttpOnly: true,
	}

	scopeCookie := &http.Cookie{
		Name:     cookiey.CookieScope,
		Value:    at.Scope,
		MaxAge:   at.ExpiresIn,
		HttpOnly: true,
	}

	isLoggedInCookie := &http.Cookie{
		Name:   cookiey.CookieIsLoggedIn,
		Value:  "true",
		MaxAge: at.ExpiresIn,
	}

	return []*http.Cookie{accessTokenCookie, refreshTokenCookie, tokenTypeCookie, scopeCookie, isLoggedInCookie}
}

type RefreshTokenRequest struct {
	RefreshToken string `validate:"required"`
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
	AccessToken string // cookie
	TokenType   string // cookie
}

func (gubt *GetUserByTokenRequest) ConstructAuthorization() string {
	return fmt.Sprintf("%s %s", gubt.TokenType, gubt.AccessToken)
}
