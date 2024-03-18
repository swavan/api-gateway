package user

import "time"

type UserStore struct {
	ID         int64  `db:"id"`
	Username   string `db:"username"`
	Password   string `db:"password"`
	FirstName  string `db:"first_name"`
	MiddleName string `db:"middle_name"`
	LastName   string `db:"last_name"`
	Email      string `db:"email"`
	Active     bool   `db:"active"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

type User struct {
	Username   string `db:"username"`
	Password   string `db:"password"`
	FirstName  string `db:"first_name"`
	MiddleName string `db:"middle_name"`
	LastName   string `db:"last_name"`
	Email      string `db:"email"`
}
