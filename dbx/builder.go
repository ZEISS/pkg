package dbx

import (
	"context"
	"database/sql"
	"errors"
	"io"

	"gorm.io/gorm"
)

// Migrator is a method that runs the migration.
type Migrator interface {
	// Migrate is a method that runs the migration.
	Migrate(context.Context, ...any) error
}

// Database provides methods for transactional operations.
type Database[R, W any] interface {
	// ReadTx starts a read only transaction.
	ReadTx(context.Context, func(context.Context, R) error) error
	// ReadWriteTx starts a read write transaction.
	ReadWriteTx(context.Context, func(context.Context, W) error) error

	Migrator
	io.Closer
}

// QueryError is an error that occurred while executing a query.
type QueryError struct {
	// Query is the query that caused the error.
	Query string
	// Err is the error that occurred.
	Err error
}

// Error implements the error interface.
func (e *QueryError) Error() string { return e.Query + ": " + e.Err.Error() }

// Unwrap implements the errors.Wrapper interface.
func (e *QueryError) Unwrap() error { return e.Err }

// NewQueryError returns a new QueryError.
func NewQueryError(query string, err error) *QueryError {
	return &QueryError{
		Query: query,
		Err:   err,
	}
}

type databaseImpl[R, W any] struct {
	r    ReadTxFactory[R]
	rw   ReadWriteTxFactory[W]
	conn *gorm.DB
}

// ReadTxFactory is a function that creates a new instance of Datastore.
type ReadTxFactory[R any] func(*gorm.DB) (R, error)

// ReadWriteTxFactory is a function that creates a new instance of Datastore.
type ReadWriteTxFactory[W any] func(*gorm.DB) (W, error)

// NewDatabase returns a new instance of db.
func NewDatabase[R, W any](conn *gorm.DB, r ReadTxFactory[R], rw ReadWriteTxFactory[W]) (Database[R, W], error) {
	return &databaseImpl[R, W]{r, rw, conn}, nil
}

// Close closes the database connection.
func (d *databaseImpl[R, W]) Close() error {
	db, err := d.conn.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

// RunMigrations runs the database migrations.
func (d *databaseImpl[R, W]) Migrate(ctx context.Context, dst ...interface{}) error {
	return d.conn.WithContext(ctx).AutoMigrate(dst...)
}

// ReadWriteTx starts a read only transaction.
func (d *databaseImpl[R, W]) ReadWriteTx(ctx context.Context, fn func(context.Context, W) error) error {
	tx := d.conn.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	rwtx, err := d.rw(tx)
	if err != nil {
		return err
	}

	err = fn(ctx, rwtx)
	if err != nil {
		tx.Rollback()
		return NewQueryError("rollback transaction", err)
	}

	if err := tx.Commit().Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
		return err
	}

	return nil
}

// ReadTx starts a read only transaction.
func (d *databaseImpl[R, W]) ReadTx(ctx context.Context, fn func(context.Context, R) error) error {
	tx := d.conn.WithContext(ctx).Begin()
	if tx.Error != nil {
		return NewQueryError("begin read transaction", tx.Error)
	}

	rtx, err := d.r(tx)
	if err != nil {
		return err
	}

	err = fn(ctx, rtx)
	if err != nil {
		tx.Rollback()
		return NewQueryError("rollback transaction", err)
	}

	if err := tx.Commit().Error; err != nil && !errors.Is(err, sql.ErrTxDone) {
		return NewQueryError("commit read transaction", err)
	}

	return nil
}
