package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	User struct {
		ID   uint64
		Name string
	}

	UserFilter struct {
		IDs   []uint64
		Names []string
	}
)

type UserRepo interface {
	CRUD[User, UserFilter]

	Search(ctx context.Context, query string) ([]User, error)
}

type Executor interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type userRepo struct {
	ext Executor
}

// NewUserRepo creates a new instance of the UserRepo type using the given
// pgxpool.Pool as the underlying database connection.
func NewUserRepo(ext Executor) UserRepo {
	return &userRepo{ext: ext}
}

// Create inserts the given user into the database and returns the newly created
// user.
func (r *userRepo) Create(ctx context.Context, in User) (User, error) {
	panic("not implemented")
}

// All returns all users in the database that match the given filter.
func (r *userRepo) All(ctx context.Context, filter UserFilter) ([]User, error) {
	panic("not implemented")
}

// One returns the first user in the database that matches the given filter.
func (r *userRepo) One(ctx context.Context, filter UserFilter) (User, error) {
	panic("not implemented")
}

// Update updates the given user in the database and returns the updated user.
func (r *userRepo) Update(ctx context.Context, in User) (User, error) {
	panic("not implemented")
}

// Delete deletes the given user from the database.
func (r *userRepo) Delete(ctx context.Context, in User) error {
	panic("not implemented")
}

// Search returns all users in the database that match the given query.
func (r *userRepo) Search(ctx context.Context, query string) ([]User, error) {
	panic("not implemented")
}
