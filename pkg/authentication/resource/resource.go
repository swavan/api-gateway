package resource

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Resource struct {
	ID          string `json:"id" db:"id"`
	Name        string `query:"name" json:"name" form:"name"`
	Description string `json:"description" db:"description"`
	Source      string `json:"source" db:"source"`
	Actions     string `json:"actions" db:"actions"`
	Modifier    string `json:"modifier" db:"modifier"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

func NewResource() *Resource {
	return &Resource{
		ID: uuid.New().String(),
	}
}

func (r *Resource) SetID(id string) *Resource {
	r.ID = id
	return r
}

func (r *Resource) SetName(name string) *Resource {
	r.Name = name
	return r
}

func (r *Resource) SetDescription(description string) *Resource {
	r.Description = description
	return r
}

func (r *Resource) SetSource(source string) *Resource {
	r.Source = source
	return r
}

func (r *Resource) SetActions(actions ...string) *Resource {
	act := strings.Split(r.Actions, ",")
	act = append(act, actions...)
	r.Actions = strings.Join(act, ",")
	return r
}

func (r *Resource) SetModifier(modifier string) *Resource {
	r.Modifier = modifier
	return r
}

func (r *Resource) SetCreatedAt(createdAt string) *Resource {
	r.CreatedAt = createdAt
	return r
}

func (r *Resource) SetUpdatedAt(updatedAt string) *Resource {
	r.UpdatedAt = updatedAt
	return r
}

func (r *Resource) GetActions() []string {
	return strings.Split(r.Actions, ",")
}

type ResourceAPI interface {
	All(ctx context.Context) ([]Resource, error)
	Find(ctx context.Context, id string) (*Resource, error)
	ActionsByID(ctx context.Context, id string) ([]string, error)
	ActionsBySource(ctx context.Context, id string) ([]string, error)
	Save(ctx context.Context, resource *Resource) error
	Delete(ctx context.Context, id string) error
	Migration(ctx context.Context) error
}

type ResourceService struct {
	database *sqlx.DB
	cfg      *Config
}

// All implements ResourceAPI.
func (r *ResourceService) All(ctx context.Context) ([]Resource, error) {
	resources := []Resource{}
	err := r.database.SelectContext(ctx, &resources, r.cfg.Scripts.FetchAll)
	return resources, err
}

func (r *ResourceService) ActionsByID(ctx context.Context, id string) ([]string, error) {
	var actions string
	err := r.database.GetContext(ctx, &actions, r.cfg.Scripts.FetchActionsByID, id)
	if err != nil {
		return nil, err
	}
	return strings.Split(actions, ","), nil
}

func (r *ResourceService) ActionsBySource(ctx context.Context, id string) ([]string, error) {
	actions := []string{}
	err := r.database.SelectContext(ctx, &actions, r.cfg.Scripts.FetchActionsBySource, id)
	if err != nil {
		return nil, err
	}
	maps := make(map[string]bool)
	for _, act := range actions {
		for _, a := range strings.Split(act, ",") {
			maps[a] = true
		}
	}
	resourceActions := []string{}
	for k := range maps {
		resourceActions = append(resourceActions, k)
	}
	return resourceActions, nil
}

// Delete implements ResourceAPI.
func (r *ResourceService) Delete(ctx context.Context, id string) error {
	_, err := r.database.ExecContext(ctx, r.cfg.Scripts.DeleteByID, id)
	return err
}

// Find implements ResourceAPI.
func (r *ResourceService) Find(ctx context.Context, id string) (*Resource, error) {
	resource := &Resource{}
	if id == "" {
		return resource, nil
	}
	err := r.database.GetContext(
		ctx,
		resource,
		r.cfg.Scripts.FetchByID,
		id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resource, nil
}

// Migration implements ResourceAPI.
func (r *ResourceService) Migration(ctx context.Context) error {
	for _, script := range r.cfg.Migration.Scripts {
		_, err := r.database.ExecContext(ctx, script)
		if err != nil {
			return err
		}
	}
	return nil
}

// Save implements ResourceAPI.
func (r *ResourceService) Save(ctx context.Context, resource *Resource) error {
	script := r.cfg.Scripts.Create
	if resource.ID != "" {
		script = r.cfg.Scripts.UpdateByID
	}

	params := []interface{}{
		resource.ID,
		resource.Name,
		resource.Description,
		resource.Source,
		strings.TrimPrefix(resource.Actions, ","),
		resource.Modifier,
	}
	_, err := r.database.ExecContext(
		ctx,
		script,
		params...,
	)

	return err
}

func New(database *sqlx.DB, cfg *Config) (ResourceAPI, error) {
	rs := &ResourceService{database, cfg.SetDefaultIfEmpty()}
	if cfg.Migration.Run {
		if err := rs.Migration(context.Background()); err != nil {
			return nil, err
		}
	}
	return rs, nil
}
