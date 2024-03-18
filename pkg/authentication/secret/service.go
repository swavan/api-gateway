package secret

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type SecretService struct {
	config   *Config
	database *sqlx.DB
}

// Delete implements SecretAPI.
func (t *SecretService) Delete(ctx context.Context, id string) error {
	_, err := t.database.ExecContext(
		ctx,
		t.config.Scripts.DeleteByID,
		id,
	)
	return err
}

// Archive implements SecretAPI.
func (t *SecretService) Archive(ctx context.Context, id string) error {
	_, err := t.database.ExecContext(
		ctx,
		t.config.Scripts.ArchiveByID,
		id,
	)
	return err
}

// Get implements SecretAPI.
func (t *SecretService) Get(ctx context.Context, id string) (*Secret, error) {
	var token *Secret
	err := t.database.
		GetContext(
			ctx,
			token,
			t.config.Scripts.FetchByDomain,
			id,
		)
	if err != nil {
		return nil, err
	}
	return token, nil

}

// GetByDomain implements SecretAPI.
func (t *SecretService) GetByDomain(ctx context.Context, domain ...string) ([]Secret, error) {
	var (
		tokens   []Secret
		filterBy = []interface{}{}
	)
	if len(domain) > 0 {
		filterBy = []interface{}{}
		for _, d := range domain {
			filterBy = append(filterBy, d)
		}
	}

	err := t.database.
		SelectContext(
			ctx,
			&tokens,
			t.config.Scripts.FetchByDomain,
			filterBy...,
		)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (t *SecretService) GetByUser(ctx context.Context, users ...string) ([]Secret, error) {
	var (
		tokens   []Secret
		filterBy = []interface{}{}
	)
	if len(users) > 0 {
		filterBy = []interface{}{}
		for _, d := range users {
			filterBy = append(filterBy, d)
		}
	}

	err := t.database.
		SelectContext(
			ctx,
			&tokens,
			t.config.Scripts.FetchByUser,
			filterBy...,
		)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Save implements SecretAPI.
func (t *SecretService) Save(ctx context.Context, content *Secret) error {
	_, err := t.database.ExecContext(
		ctx,
		t.config.Scripts.Save,
		content.ID,
		content.Description,
		content.Type,
		content.Domain,
		content.IssueAt,
		content.ExpiresAt,
		content.AlertTo,
		content.Modifier,
	)
	return err
}

// Migration implements SecretAPI.
func (t *SecretService) Migration(context.Context) error {
	for _, script := range t.config.Migration.Scripts {
		if _, err := t.database.ExecContext(context.Background(), script); err != nil {
			return err
		}
	}
	return nil

}

func New(database *sqlx.DB, config *Config) (SecretAPI, error) {
	ts := &SecretService{
		config:   config.SetDefaultIfEmpty(),
		database: database,
	}
	if config.Migration.Run {
		if err := ts.Migration(context.Background()); err != nil {
			return nil, err
		}
	}
	return ts, nil
}
