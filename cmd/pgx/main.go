package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/khanfromasia/datastore/internal/storage"
)

func main() {
	ctx := context.TODO()

	ds, err := storage.NewConnectedPostgresDB("putYourURLHere")
	if err != nil {
		panic(err)
	}

	myUserID := 1
	productName := "book"

	err = ds.InTransaction(ctx, func(ctx context.Context, txStore storage.DataStore) error {
		user, err := txStore.User().One(ctx, storage.UserFilter{IDs: []uint64{uint64(myUserID)}})
		if err != nil {
			return err
		}

		productName = fmt.Sprintf("%s for %s", productName, user.Name)

		return nil
	}, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		panic(err)
	}
}
