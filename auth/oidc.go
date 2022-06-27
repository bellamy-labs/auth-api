package auth

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type Provider struct {
	AuthHandler     echo.HandlerFunc
	CallbackHandler echo.HandlerFunc
}

type Service struct {
	Issuer       string
	Provider     string
	ClientID     string
	ClientSecret string
}

var Verifier *oidc.IDTokenVerifier

func (s *Service) InitService() (*oauth2.Config, error) {
	provider, err := oidc.NewProvider(context.Background(), s.Issuer)
	if err != nil {
		return &oauth2.Config{}, err
	}

	oidcConfig := &oidc.Config{
		ClientID: s.ClientID,
	}
	Verifier = provider.Verifier(oidcConfig)

	return &oauth2.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8080/api/v1/auth/" + s.Provider + "/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}, nil
}
