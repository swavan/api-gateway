package authentication

import (
	"os"

	"github.com/swavan.io/gateway/pkg/authentication/access"
	"github.com/swavan.io/gateway/pkg/authentication/domain"
	"github.com/swavan.io/gateway/pkg/authentication/key"
	"github.com/swavan.io/gateway/pkg/authentication/oidc"
	"github.com/swavan.io/gateway/pkg/authentication/resource"
	"github.com/swavan.io/gateway/pkg/authentication/role"
	"github.com/swavan.io/gateway/pkg/authentication/secret"
	"github.com/swavan.io/gateway/pkg/authentication/user"
)

type AuthConfig struct {
	Confidential   string               `mapstructure:"confidential"`
	Migration      bool                 `mapstructure:"migration"`
	OpenIDConnects []oidc.OpenIDConnect `mapstructure:"oidc"`
	AccessConfig   access.Config        `mapstructure:"access"`
	KeyConfig      key.Config           `mapstructure:"key"`
	DomainConfig   domain.Config        `mapstructure:"domain"`
	RoleConfig     role.Config          `mapstructure:"role"`
	UserConfig     user.Config          `mapstructure:"user"`
	ResourceConfig resource.Config      `mapstructure:"resource"`
	SecretConfig   secret.Config        `mapstructure:"secret"`
	IgnoreAccess   []string             `mapstructure:"ignore_access"`
	SuperAdmins    []struct {
		Domain   string   `mapstructure:"domain"`
		Resource string   `mapstructure:"resource"`
		Role     string   `mapstructure:"role"`
		action   string   `mapstructure:"actions"`
		User     []string `mapstructure:"users"`
	} `mapstructure:"admins"`
	Users []struct {
		Username  string `mapstructure:"username"`
		Password  string `mapstructure:"password"`
		Email     string `mapstructure:"email"`
		FirstName string `mapstructure:"first_name"`
		LastName  string `mapstructure:"last_name"`
		NoneUser  bool   `mapstructure:"none_user"`
	} `mapstructure:"users"`
}

func NewConfig() *AuthConfig {
	return new(AuthConfig)
}

func (c *AuthConfig) Name() string {
	return os.Getenv("AUTH_CONFIG_FILE_NAME")
}

func (c *AuthConfig) Path() string {
	return os.Getenv("AUTH_CONFIG_FILE_LOCATION")
}

func (c *AuthConfig) Extension() string {
	return os.Getenv("AUTH_CONFIG_FILE_EXTENSION")
}
