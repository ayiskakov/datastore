package storage

import (
	"context"
	"io"

	"github.com/jackc/pgx/v5"
)

type TxFunc[T any] func(ctx context.Context, txStore T) error

type Transactor[T any] interface {
	InTransaction(ctx context.Context, txFunc TxFunc[T], opts pgx.TxOptions) error
}

type DataStore interface {
	Transactor[DataStore]
	io.Closer

	User() UserRepo
	//Product() ProductRepo
}

type CRUD[Entity, Filter any] interface {
	Create(ctx context.Context, in Entity) (Entity, error)
	All(ctx context.Context, filter Filter) ([]Entity, error)
	One(ctx context.Context, filter Filter) (Entity, error)
	Update(ctx context.Context, in Entity) (Entity, error)
	Delete(ctx context.Context, in Entity) error
}
