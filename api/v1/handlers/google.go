package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bellamy-labs/auth-api/auth"
	"github.com/bellamy-labs/auth-api/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var (
	googleIssuer            = "https://accounts.google.com"
	googleProviderName      = "google"
	googleStateCookieName   = "google_state"
	googleNonceCookieName   = "google_nonce"
	googleStateCookieExpiry = time.Now().UTC().Add(time.Hour)
	googleNonceCookieExpiry = time.Now().UTC().Add(time.Hour)
	//	googleCookiePath        = "/api/v1/auth/google"
)

// google oauth config
func InitGoogle() (*auth.Provider, error) {
	s := auth.Service{
		Issuer:       googleIssuer,
		Provider:     googleProviderName,
		ClientID:     config.Cfg.GoogleClientID,
		ClientSecret: config.Cfg.GoogleClientSecret,
	}

	googleCfg, err := s.InitService()
	if err != nil {
		return &auth.Provider{}, err
	}

	// make authenticator handler
	googleAuth := func(c echo.Context) error {
		state := randString(32)
		nonce := randString(32)

		sc := &http.Cookie{
			Path:     "/",
			Name:     googleStateCookieName,
			Value:    state,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Expires:  googleStateCookieExpiry,
		}

		nc := &http.Cookie{
			Path:     "/",
			Name:     googleNonceCookieName,
			Value:    nonce,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Expires:  googleNonceCookieExpiry,
		}

		c.SetCookie(sc)
		c.SetCookie(nc)

		return c.Redirect(http.StatusFound, googleCfg.AuthCodeURL(state, oidc.Nonce(nonce)))
	}

	// make callback handler
	googleCallback := func(c echo.Context) error {
		state, err := c.Cookie(googleStateCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return &ErrorResponse{Status: http.StatusUnauthorized,
					Message: http.StatusText(http.StatusUnauthorized)}
			}

			return err
		}

		if c.Request().URL.Query().Get("state") != state.Value {
			return &ErrorResponse{Status: http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized)}
		}

		oauth2Token, err := googleCfg.Exchange(context.Background(), c.Request().URL.Query().Get("code"))
		if err != nil {
			return &ErrorResponse{Status: http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized)}
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			return &ErrorResponse{Status: http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized)}
		}
		idToken, err := auth.Verifier.Verify(context.Background(), rawIDToken)
		if err != nil {
			return &ErrorResponse{Status: http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized)}
		}

		nonce, err := c.Cookie(googleNonceCookieName)
		if err != nil {
			return &ErrorResponse{Status: http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized)}
		}
		if idToken.Nonce != nonce.Value {
			return &ErrorResponse{Status: http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized)}
		}

		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		}{oauth2Token, new(json.RawMessage)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, resp)
	}

	p := auth.Provider{
		AuthHandler:     googleAuth,
		CallbackHandler: googleCallback,
	}

	return &p, nil
}
