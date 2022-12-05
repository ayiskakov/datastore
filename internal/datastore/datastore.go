package datastore

import (
	"context"
	"io"
)

type TxFunc[T any] func(txStore T) error

type Transactor[T any] interface {
	InTransaction(txFunc TxFunc[T]) error
}

type DataStore interface {
	Transactor[DataStore]
	io.Closer

	User() UserRepo
	Product() ProductRepo
}

type CRUD[Entity, Filter any] interface {
	Create(ctx context.Context, in Entity) (Entity, error)
	All(ctx context.Context, filter Filter) ([]Entity, error)
	One(ctx context.Context, filter Filter) (Entity, error)
	Update(ctx context.Context, in Entity) (Entity, error)
	Delete(ctx context.Context, in Entity) error
}
