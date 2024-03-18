package authentication

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/swavan.io/gateway/pkg/authentication/access"
	"github.com/swavan.io/gateway/pkg/authentication/domain"
	"github.com/swavan.io/gateway/pkg/authentication/key"
	"github.com/swavan.io/gateway/pkg/authentication/oidc"
	"github.com/swavan.io/gateway/pkg/authentication/resource"
	"github.com/swavan.io/gateway/pkg/authentication/role"
	"github.com/swavan.io/gateway/pkg/authentication/secret"
	"github.com/swavan.io/gateway/pkg/authentication/user"

	"github.com/google/uuid"
)

type AuthenticationAPI interface {
	Key() key.KeyManagerAPI
	Role() role.RoleAPI
	Resource() resource.ResourceAPI
	Domain() domain.DomainAPI
	User() user.UserAPI
	Secret() secret.SecretAPI
	Access() access.API
	OIDC() oidc.OauthClients
	Config() *AuthConfig
}

type Authentication struct {
	resource resource.ResourceAPI
	role     role.RoleAPI
	domain   domain.DomainAPI
	user     user.UserAPI
	access   access.API
	secret   secret.SecretAPI
	oidc     oidc.OauthClients
	cfg      *AuthConfig
	key      key.KeyManagerAPI
}

// Token implements AuthenticationAPI.
func (a *Authentication) Secret() secret.SecretAPI {
	return a.secret
}

// AccessControl implements AuthenticationAPI.
func (a *Authentication) Access() access.API {
	return a.access
}

// Resource implements AuthenticationAPI.
func (a *Authentication) Resource() resource.ResourceAPI {
	return a.resource
}

// Domain implements AuthenticationAPI.
func (a *Authentication) Domain() domain.DomainAPI {
	return a.domain
}

// Role implements AuthenticationAPI.
func (a *Authentication) Role() role.RoleAPI {
	return a.role
}

// User implements AuthenticationAPI.
func (a *Authentication) User() user.UserAPI {
	return a.user
}

// OIDC implements AuthenticationAPI.
func (a *Authentication) OIDC() oidc.OauthClients {
	return a.oidc
}

// Config implements AuthenticationAPI.
func (a *Authentication) Config() *AuthConfig {
	return a.cfg
}

// Key implements AuthenticationAPI.
func (a *Authentication) Key() key.KeyManagerAPI {
	return a.key
}

func New(dep *sqlx.DB, configs ...*AuthConfig) (AuthenticationAPI, error) {
	cfg := new(AuthConfig)
	for _, c := range configs {
		cfg = c
	}
	usr, err := user.New(dep, &cfg.UserConfig)
	if err != nil {
		return nil, err
	}

	rl, err := role.New(dep, &cfg.RoleConfig)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(dep, &cfg.ResourceConfig)
	if err != nil {
		return nil, err
	}

	dm, err := domain.New(dep, &cfg.DomainConfig)
	if err != nil {
		return nil, err
	}

	key, err := key.New(dep, &cfg.KeyConfig, os.Getenv(cfg.Confidential))
	if err != nil {
		return nil, err
	}

	access, err := access.New(cfg.AccessConfig.Policy, dep, cfg.Migration)
	if err != nil {
		return nil, err
	}

	oidc, err := oidc.New(context.Background(), cfg.OpenIDConnects)
	if err != nil {
		return nil, err
	}

	sec, err := secret.New(dep, &cfg.SecretConfig)
	if err != nil {
		return nil, err
	}

	if err := CreateUsers(usr, cfg); err != nil {
		return nil, err
	}

	if err := access.Enforcer().LoadPolicy(); err != nil {
		return nil, err
	}

	auth := &Authentication{
		user:     usr,
		resource: res,
		role:     rl,
		domain:   dm,
		access:   access,
		oidc:     oidc,
		cfg:      cfg,
		key:      key,
		secret:   sec,
	}

	return auth, nil
}

func CreateUsers(usr user.UserAPI, cfg *AuthConfig) error {
	for _, u := range cfg.Users {
		newUser := user.NewUser().
			SetID(uuid.New().String()).
			SetEmail(u.Email).
			SetUsername(u.Username).
			SetGivenName(u.FirstName).
			SetFamilyName(u.LastName).
			SetEmailVerified(false)
		if err := usr.Save(context.Background(), newUser); err != nil {
			return err
		}

		if err := usr.ChangePassword(context.Background(), newUser.Username, u.Password); err != nil {
			return err
		}
	}
	return nil
}

func (auth Authentication) SetupUserToDomains(users []string, dom string, role string) error {
	roleTemplate := "role:%s"
	for _, username := range users {

		_, err := auth.access.Enforcer().AddRoleForUserInDomain(
			username,
			fmt.Sprintf(roleTemplate, role),
			dom,
		)
		if err != nil {
			return err
		}

		acc, er := auth.user.FindByUsername(context.Background(), username)
		if er != nil {
			return err
		}

		if err := auth.user.Save(context.Background(), acc.SetUsername(username).
			SetDomains(dom)); err != nil {
			return err
		}
	}

	return nil
}

func (auth Authentication) SetupSuperUser(usr user.UserAPI, cfg *AuthConfig) error {
	roleTemplate := "role:%s"

	for _, profile := range cfg.SuperAdmins {

		dom, err := auth.domain.FetchByName(
			context.Background(),
			profile.Domain,
		)
		if err != nil {
			return err
		}

		if dom == nil {
			dom = domain.NewDomain().SetName(profile.Domain).NewID()
			if err := auth.domain.Create(context.Background(), dom); err != nil {
				return err
			}
		}

		if err := auth.role.Save(
			context.Background(),
			role.NewRole().SetName(profile.Role)); err != nil {
			return err
		}

		status := auth.access.Enforcer().HasPolicy(
			fmt.Sprintf(roleTemplate, profile.Role), dom.ID, profile.Resource, "*")

		if !status {
			if _, err := auth.access.Enforcer().AddPolicy(
				fmt.Sprintf(roleTemplate, profile.Role), // Role
				dom.ID, profile.Resource, "ANY"); err != nil {
				return err
			}
		}

		if err := auth.SetupUserToDomains(profile.User, dom.ID, profile.Role); err != nil {
			return err
		}

	}

	return nil

}
