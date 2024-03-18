package oidc

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OpenIDConnect struct {
	ID             string   `mapstructure:"id"`
	Name           string   `mapstructure:"name"`
	Icon           string   `mapstructure:"icon"`
	ClientID       string   `mapstructure:"client"`
	Secret         string   `mapstructure:"secret"`
	RedirectURL    string   `mapstructure:"redirect"`
	Scopes         []string `mapstructure:"scopes"`
	IssuerEndpoint string   `mapstructure:"issuer"`
	Enabled        bool     `mapstructure:"enabled"`
}

type OauthClients map[string]*OauthClient

type AuthConfig struct {
	OauthClients OauthClients
}

type OauthClient struct {
	ID           string
	Name         string
	Logout       string
	ClientID     string
	ClientSecret string
	PublicKey    string
	Provider     *oidc.Provider
	AuthConfig   oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

func New(ctx context.Context, clientConfigs []OpenIDConnect) (OauthClients, error) {
	clients := make(OauthClients)
	for _, cfg := range clientConfigs {
		if !cfg.Enabled {
			continue
		}
		client, err := newOauthClient(ctx, cfg)
		if err != nil {
			return nil, err
		}
		clients[cfg.ID] = client
	}
	return clients, nil
}

func newOauthClient(ctx context.Context, cfg OpenIDConnect) (*OauthClient, error) {
	provider, err := newProvider(ctx, cfg.IssuerEndpoint)
	if err != nil {
		return nil, err
	}
	return &OauthClient{
		ID:           cfg.ID,
		Name:         cfg.Name,
		Provider:     provider,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.Secret,
		AuthConfig:   newAuthConfig(cfg, provider),
		Verifier:     newVerifier(provider, cfg.ClientID),
	}, nil
}

func newProvider(ctx context.Context, issuerURL string) (*oidc.Provider, error) {
	return oidc.NewProvider(ctx, issuerURL)
}

func newAuthConfig(cfg OpenIDConnect, provider *oidc.Provider) oauth2.Config {
	return oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.Secret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
	}
}

func newVerifier(provider *oidc.Provider, clientId string) *oidc.IDTokenVerifier {
	return provider.Verifier(&oidc.Config{
		ClientID: clientId,
	})
}
