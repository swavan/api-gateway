package secret

import "context"

type SecretAPI interface {
	Migration(context.Context) error
	Save(context.Context, *Secret) error
	Get(context.Context, string) (*Secret, error)
	Delete(context.Context, string) error
	Archive(context.Context, string) error
	GetByUser(context.Context, ...string) ([]Secret, error)
	GetByDomain(context.Context, ...string) ([]Secret, error)
}
