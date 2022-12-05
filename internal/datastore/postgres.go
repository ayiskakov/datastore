package datastore

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // custom postgres driver
	"github.com/jmoiron/sqlx"
)

var _ DataStore = (*PostgresDB)(nil)

const (
	connMaxLifetime     = time.Minute * 10
	connMaxIdleLifetime = time.Second * 30
	maxOpenConns        = 100
	maxIdleConns        = maxOpenConns / 3
)

type PostgresDB struct {
	executor sqlx.ExtContext

	user    UserRepo
	product ProductRepo
}

func NewConnectedPostgresDB(dbURL string) (PostgresDB, error) {
	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		return PostgresDB{}, fmt.Errorf("failed to connect to db: %w", err)
	}

	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	return NewPostgresDB(db), nil
}

func NewPostgresDB(executor sqlx.ExtContext) PostgresDB {
	ds := PostgresDB{
		executor: executor,
	}
	ds.setup()

	return ds
}

func (d *PostgresDB) setup() {
	//d.user = NewUserRepo(d.executor)       // TODO Implement constructor
	//d.product = NewProductRepo(d.executor) // TODO Implement constructor
}

func (d PostgresDB) InTransaction(txFunc TxFunc[DataStore]) error {
	db, isDB := d.executor.(*sqlx.DB)
	if !isDB {
		return errors.New("nested transactions are not allowed")
	}

	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		//  defer tx.Rollback() is safe even if tx.Commit() will be called first in a non-error condition
		_ = tx.Rollback()
	}()

	txStore := NewPostgresDB(tx)
	if err = txFunc(txStore); err != nil {
		return fmt.Errorf("failed to exec tx: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

func (d PostgresDB) Close() error {
	db, isDB := d.executor.(*sqlx.DB)
	if isDB {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close db: %w", err)
		}

		return nil
	}

	return nil
}

func (d PostgresDB) User() UserRepo {
	return d.user
}

func (d PostgresDB) Product() ProductRepo {
	return d.product
}
