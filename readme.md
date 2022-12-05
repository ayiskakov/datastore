# Data Store

Basing on [gist of RA9dev](https://gist.github.com/ra9dev/1ae6fba7382aaa23f42e34c2a0164f9d) I've created a simple data store for my projects using v5/pgxpool

This package provides a simple data store that supports CRUD operations for a user and product. It uses a database connection pool to manage the connections to the database and provides transaction management using pgx.TxOptions.

Usage

To create a new data store, call the NewConnectedPostgresDB function with a valid database URL. This will create a new connection pool and connect to the database using the provided URL.
```
db, err := storage.NewConnectedPostgresDB("putYourURLHere")
if err != nil {
    panic(err)
}
```
The data store provides methods for working with users and products. To create a new user, call the User method to get a user repository and use the Create method to create a new user.

```
user := storage.User{
    Name: "John Doe",
}
user, err := db.User().Create(context.TODO(), user)
if err != nil {
    panic(err)
}
```
The data store also provides a InTransaction method that allows you to run a function in the context of a transaction. If the function returns an error, the transaction is rolled back. Otherwise, the transaction is committed.

```
err := db.InTransaction(context.TODO(), func(ctx context.Context, txStore storage.DataStore) error {
    user, err := txStore.User().One(ctx, storage.UserFilter{IDs: []uint64{1}})
    if err != nil {
        return err
    }


    product, err := txStore.Product().Create(ctx, storage.Product{
        Name: "Book for John Doe",
    })
    if err != nil {
        return err
    }
    return nil
}, pgx.TxOptions{})
if err != nil {
    panic(err)
}
```
When you are done with the data store, be sure to call the Close method to close the underlying database connections.