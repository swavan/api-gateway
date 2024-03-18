package role

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Role struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Modifier    string `json:"modifier" db:"modifier"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

func NewRole() *Role {
	return &Role{
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

func (r *Role) SetID(id int64) *Role {
	r.ID = id
	return r
}

func (r *Role) SetName(name string) *Role {
	r.Name = name
	return r
}

func (r *Role) SetDescription(description string) *Role {
	r.Description = description
	return r
}

func (r *Role) SetModifier(modifier string) *Role {
	r.Modifier = modifier
	return r
}

type RoleAPI interface {
	All(ctx context.Context) ([]Role, error)
	Find(ctx context.Context, id int64) (*Role, error)
	Save(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id int64) error
	Migration(ctx context.Context) error
}

type RoleService struct {
	database *sqlx.DB
	cfg      *Config
}

// Migrate implements RoleAPI.
func (rs *RoleService) Migration(ctx context.Context) error {
	if !rs.cfg.Migration.Run {
		return nil
	}
	for _, script := range rs.cfg.Migration.Scripts {
		_, err := rs.database.ExecContext(ctx, script)
		if err != nil {
			return err
		}
	}

	return nil
}

// All implements RoleAPI.
func (rs *RoleService) All(ctx context.Context) ([]Role, error) {
	roles := []Role{}
	err := rs.database.SelectContext(
		ctx,
		&roles,
		rs.cfg.Scripts.FetchAll)
	return roles, err
}

// Save implements RoleAPI.
func (rs *RoleService) Save(ctx context.Context, role *Role) error {
	if role.ID == 0 {
		_, err := rs.database.ExecContext(
			ctx,
			rs.cfg.Scripts.Save,
			role.Name,
			role.Description,
			role.Modifier)
		return err
	}
	_, err := rs.database.ExecContext(
		ctx,
		rs.cfg.Scripts.UpdateByID,
		role.ID,
		role.Name,
		role.Description,
		role.Modifier)
	return err
}

// Delete implements RoleAPI.
func (rs *RoleService) Delete(ctx context.Context, id int64) error {
	_, err := rs.database.ExecContext(
		ctx,
		rs.cfg.Scripts.DeleteByID,
		id)
	return err
}

// Find implements RoleAPI.
func (rs *RoleService) Find(ctx context.Context, id int64) (*Role, error) {
	role := new(Role)
	err := rs.database.GetContext(
		ctx,
		role,
		rs.cfg.Scripts.FetchByID,
		id)
	if err != nil && err == sql.ErrNoRows {
		return role, nil
	}
	return role, err
}

func New(database *sqlx.DB, cfg *Config) (RoleAPI, error) {
	rl := &RoleService{
		database: database,
		cfg:      cfg.SetDefaultIfEmpty(),
	}
	return rl, rl.Migration(context.Background())
}
