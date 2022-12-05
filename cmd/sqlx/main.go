package main

import (
	"context"
	"fmt"

	"github.com/khanfromasia/datastore/internal/datastore"
)

func main() {
	ctx := context.TODO()

	ds, err := datastore.NewConnectedPostgresDB("putYourURLHere")
	if err != nil {
		panic(err)
	}

	myUserID := 1
	productName := "book"

	err = ds.InTransaction(func(txStore datastore.DataStore) error {
		user, err := txStore.User().One(ctx, datastore.UserFilter{IDs: []uint64{uint64(myUserID)}})
		if err != nil {
			return err
		}

		productName = fmt.Sprintf("%s for %s", productName, user.Name)

		return nil
	})

	if err != nil {
		panic(err)
	}
}
