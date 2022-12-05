package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connMaxLifetime     = time.Minute * 10
	connMaxIdleLifetime = time.Second * 30
	maxOpenConns        = 100
	maxIdleConns        = maxOpenConns / 3
)

type PostgresDB struct {
	pool     *pgxpool.Pool
	executor Executor
	user     UserRepo
	product  ProductRepo
}

func NewConnectedPostgresDB(dbURL string) (PostgresDB, error) {
	poolConfig, err := pgxpool.ParseConfig(dbURL)

	poolConfig.MaxConnIdleTime = connMaxIdleLifetime
	poolConfig.MaxConnLifetime = connMaxLifetime
	poolConfig.MaxConns = maxOpenConns
	poolConfig.MinConns = maxIdleConns

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return PostgresDB{}, fmt.Errorf("failed to connect to db: %w", err)
	}

	return NewPostgresDB(pool, pool), nil
}

func NewPostgresDB(pool *pgxpool.Pool, ext Executor) PostgresDB {
	ds := PostgresDB{
		pool:     pool,
		executor: ext,
	}
	ds.setup()

	return ds
}

func (d *PostgresDB) setup() {
	d.user = NewUserRepo(d.executor)
	//d.product = NewProductRepo(d.pool)
}

func (d PostgresDB) User() UserRepo {
	return d.user
}

func (d PostgresDB) Product() ProductRepo {
	return d.product
}

// Close closes the underlying database connection.
func (d PostgresDB) Close() error {
	d.pool.Close()
	return nil
}

// InTransaction starts a new transaction and runs the given function in the
// context of the transaction. If the function returns an error, the transaction
// is rolled back. Otherwise, the transaction is committed.
func (d PostgresDB) InTransaction(ctx context.Context, txFunc TxFunc[DataStore], opts pgx.TxOptions) error {
	tx, err := d.pool.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer rolling back the transaction in case of an error.
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Call the given function in the context of the transaction.
	if txFunc != nil {
		txStore := NewPostgresDB(d.pool, tx)
		if err := txFunc(ctx, txStore); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
