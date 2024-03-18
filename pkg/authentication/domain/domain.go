package domain

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/jmoiron/sqlx"
)

type Domain struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Modifier    string `json:"modifier" db:"modifier"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

func NewDomain() *Domain {
	return &Domain{
		ID:        "",
		Name:      "",
		Modifier:  "",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

func (d *Domain) NewID() *Domain {
	d.ID = string(base64.RawURLEncoding.EncodeToString([]byte(d.Name)))
	return d
}

func (d *Domain) SetID(id string) *Domain {
	d.ID = id
	return d
}

func (d *Domain) SetName(name string) *Domain {
	d.Name = name
	return d
}

func (d *Domain) SetDescription(description string) *Domain {
	d.Description = description
	return d
}

func (d *Domain) SetModifier(modifier string) *Domain {
	d.Modifier = modifier
	return d
}

type DomainAPI interface {
	All(ctx context.Context, ids ...string) ([]Domain, error)
	Find(ctx context.Context, id string) (*Domain, error)
	FetchByName(ctx context.Context, domainName string) (*Domain, error)
	Create(ctx context.Context, domain *Domain) error
	Save(ctx context.Context, domain *Domain) error
	Update(ctx context.Context, domain *Domain) error
	Delete(ctx context.Context, id string) error
	Migration(ctx context.Context) error
}

type DomainService struct {
	database *sqlx.DB
	config   *Config
}

// Migrate implements DomainAPI.
func (ds *DomainService) Migration(ctx context.Context) error {
	if !ds.config.Migration.Run {
		return nil
	}
	for _, script := range ds.config.Migration.Scripts {
		_, err := ds.database.ExecContext(ctx, script)
		if err != nil {
			return err
		}
	}
	return nil
}

// All implements DomainAPI.
func (ds *DomainService) All(ctx context.Context, ids ...string) ([]Domain, error) {
	domains := []Domain{}
	if len(ids) > 0 {
		query, args, err := sqlx.In(ds.config.Scripts.FetchByIDs, ids)
		if err != nil {
			return domains, err
		}
		// query = ds.database.Rebind(query)
		// if err := ds.database.SelectContext(ctx, &domains, query, args...); err != nil {
		// 	return domains, err
		// }
		err = ds.database.SelectContext(
			ctx,
			&domains,
			ds.database.Rebind(query),
			args...,
		)
		return domains, err
	}
	err := ds.database.SelectContext(
		ctx,
		&domains,
		ds.config.Scripts.FetchAll)
	return domains, err
}

func (ds *DomainService) FetchByName(ctx context.Context, domainName string) (*Domain, error) {
	dm := NewDomain()
	err := ds.database.
		Get(
			dm,
			ds.config.Scripts.FetchByName,
			domainName,
		)
	if err != nil {
		return nil, nil
	}
	return dm, err
}

// Create implements DomainAPI.
func (ds *DomainService) Create(ctx context.Context, domain *Domain) error {

	_, err := ds.database.
		ExecContext(
			ctx,
			ds.config.Scripts.Save,
			domain.ID,
			domain.Name,
			domain.Description,
			domain.Modifier,
		)
	return err
}

func (ds *DomainService) Save(ctx context.Context, domain *Domain) error {
	if domain.ID == "" {
		return ds.Create(ctx, domain.NewID())
	}
	return ds.Update(ctx, domain)
}

// Find implements DomainAPI.
func (ds *DomainService) Find(ctx context.Context, id string) (*Domain, error) {
	dm := NewDomain()
	err := ds.database.
		Get(
			dm,
			ds.config.Scripts.FetchByID,
			id,
		)
	if err != nil && err == sql.ErrNoRows {
		return dm, nil
	}
	return dm, err
}

// Update implements DomainAPI.
func (ds *DomainService) Update(ctx context.Context, domain *Domain) error {
	_, err := ds.database.
		ExecContext(
			ctx,
			ds.config.Scripts.UpdateByID,
			domain.Name,
			domain.Description,
			domain.Modifier,
			domain.ID,
		)
	return err
}

// Delete implements DomainAPI.
func (ds *DomainService) Delete(ctx context.Context, id string) error {
	_, err := ds.database.ExecContext(ctx, ds.config.Scripts.DeleteByID, id)
	return err
}

func New(database *sqlx.DB, cfg *Config) (DomainAPI, error) {
	dm := &DomainService{
		database: database,
		config:   cfg.SetDefaultIfEmpty(),
	}
	if err := dm.Migration(context.Background()); err != nil {
		return nil, err
	}
	return dm, nil
}
