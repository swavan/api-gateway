package user

import (
	"context"
	"database/sql"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserAPI interface {
	Migration(ctx context.Context) error
	All(ctx context.Context) ([]User, error)
	Find(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	AddDomains(ctx context.Context, username string, domains ...string) error
	RemoveDomains(ctx context.Context, username string, domains ...string) error
	ChangePassword(ctx context.Context, username string, password string) error
	GetUserForCredential(ctx context.Context, username string) (*UserStore, error)
	GetDomains(ctx context.Context, username string) ([]string, error)
	Save(ctx context.Context, user *User) error
}

type User struct {
	ID                string `json:"id" db:"id"`
	Username          string `json:"username" db:"user_name"`
	PreferredUsername string `json:"preferred_username" db:"preferred_username"`
	Name              string `json:"name" db:"name"`
	GivenName         string `json:"given_name" db:"given_name"`
	FamilyName        string `json:"family_name" db:"family_name"`
	Email             string `json:"email" db:"email"`
	EmailVerified     bool   `json:"email_verified" db:"email_verified"`
	Avatar            string `json:"avatar" db:"avatar"`
	Domains           string `json:"domains" db:"domains"`
	NoneUser          bool   `json:"none_user" db:"none_user"`
	CreatedAt         string `json:"created_at" db:"created_at"`
}

type UserStore struct {
	Password string `db:"secret"`
	*User
}

func NewUser() *User {
	return &User{}
}

func (u *User) IsNew() bool {
	return u.ID == ""
}

func (u *User) SetID(id string) *User {
	u.ID = id
	return u
}

func (u *User) SetUsername(username string) *User {
	u.Username = username
	return u
}

func (u *UserStore) SetPassword(password string) *UserStore {
	u.Password = password
	return u
}

func (u *User) SetPreferredUsername(preferredUsername string) *User {
	u.PreferredUsername = preferredUsername
	return u
}

func (u *User) SetGivenName(givenName string) *User {
	u.GivenName = givenName
	return u
}

func (u *User) SetName(name string) *User {
	u.Name = name
	return u
}

func (u *User) SetFamilyName(familyName string) *User {
	u.FamilyName = familyName
	return u
}

func (u *User) SetEmailVerified(emailVerified bool) *User {
	u.EmailVerified = emailVerified
	return u
}

func (u *User) SetAvatar(avatar string) *User {
	u.Avatar = avatar
	return u
}

func (u *User) SetEmail(email string) *User {
	u.Email = email
	return u
}

func (u *User) SetDomains(domains ...string) *User {
	doms := u.GetDomains()
	for _, domain := range domains {
		if !slices.Contains(doms, domain) {
			doms = append(doms, domain)
		}
	}
	u.Domains = strings.TrimPrefix(strings.Join(doms, ","), ",")
	return u
}

func (u *User) GetDomains() []string {
	return strings.Split(u.Domains, ",")
}

type UserService struct {
	database *sqlx.DB
	cfg      *Config
}

// AddDomains implements UserAPI.
func (us *UserService) AddDomains(ctx context.Context, username string, domains ...string) error {
	old, err := us.GetDomains(ctx, username)
	if err != nil {
		return err
	}
	return us.saveDomains(ctx, username, append(old, domains...))
}

func (us *UserService) ChangePassword(ctx context.Context, username string, password string) error {
	_, err := us.database.
		ExecContext(
			ctx,
			us.cfg.Scripts.ChangePassword,
			password,
			username,
		)
	return err
}

func (us *UserService) saveDomains(ctx context.Context, username string, domains []string) error {
	_, err := us.database.
		ExecContext(
			ctx,
			us.cfg.Scripts.UpdateByUsername,
			strings.Join(domains, ","),
			username,
		)
	return err
}

// GetDomains implements UserAPI.
func (us *UserService) GetDomains(ctx context.Context, username string) ([]string, error) {
	var domains string
	err := us.database.GetContext(
		ctx,
		&domains,
		us.cfg.Scripts.FetchDomainByUsername,
		username)
	if err != nil && err != sql.ErrNoRows {
		return []string{}, err
	}
	return strings.Split(domains, ","), err
}

// RemoveDomains implements UserAPI.
func (us *UserService) RemoveDomains(ctx context.Context, username string, domains ...string) error {
	old, err := us.GetDomains(ctx, username)
	if err != nil {
		return err
	}
	newDomains := []string{}
	for _, domain := range old {
		if !slices.Contains(domains, domain) {
			newDomains = append(newDomains, domain)
		}
	}
	return us.saveDomains(ctx, username, newDomains)
}

func New(dep *sqlx.DB, cfg *Config) (UserAPI, error) {
	us := &UserService{
		database: dep,
		cfg:      cfg.SetDefaultIfEmpty(),
	}
	if err := us.Migration(context.Background()); err != nil {
		return nil, err
	}

	return us, nil
}

func (us *UserService) Migration(ctx context.Context) error {
	if !us.cfg.Migration.Run {
		return nil
	}
	for _, script := range us.cfg.Migration.Scripts {
		_, err := us.database.ExecContext(ctx, script)
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *UserService) All(ctx context.Context) ([]User, error) {
	users := []User{}
	err := us.database.
		SelectContext(
			ctx,
			&users,
			us.cfg.Scripts.FetchAll)
	return users, err
}

func (us *UserService) Find(ctx context.Context, id string) (*User, error) {
	user := NewUser()
	err := us.database.
		GetContext(
			ctx,
			user,
			us.cfg.Scripts.FetchByID,
			id)
	return user, err
}

func (us *UserService) FindByUsername(ctx context.Context, username string) (*User, error) {
	user := NewUser()
	err := us.database.
		GetContext(
			ctx,
			user,
			us.cfg.Scripts.FetchByUsername,
			username)
	if err != nil && err == sql.ErrNoRows {
		return user, nil
	}
	return user, err
}

func (us *UserService) GetUserForCredential(ctx context.Context, username string) (*UserStore, error) {
	user := &UserStore{
		User: NewUser(),
	}
	err := us.database.
		GetContext(
			ctx,
			user,
			us.cfg.Scripts.CheckCredentials,
			username)
	return user, err
}

func (us *UserService) Save(ctx context.Context, user *User) error {
	_, err := us.database.ExecContext(ctx, us.cfg.Scripts.Save,
		uuid.New().String(),
		user.Username,
		user.PreferredUsername,
		user.Name,
		user.GivenName,
		user.FamilyName,
		user.Username,
		user.EmailVerified,
		user.Avatar,
		user.Domains,
		user.NoneUser,
	)
	return err
}

func (us *UserService) Delete(ctx context.Context, userName string) error {
	_, err := us.database.
		ExecContext(ctx, us.cfg.Scripts.DeleteByID,
			userName)
	return err
}
