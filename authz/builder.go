package authz

import (
	"context"

	"github.com/openfga/go-sdk/client"
)

// Store is an interface that provides methods for transactional operations on the authz database.
type Store[Tx any] interface {
	// WriteTx starts a read write transaction.
	WriteTx(context.Context, func(context.Context, Tx) error) error
}

// AuthzError is an error that occurred while executing a query.
type AuthzError struct {
	// Op is the operation that caused the error.
	Op string
	// Err is the error that occurred.
	Err error
}

// Error implements the error interface.
func (e *AuthzError) Error() string { return e.Op + ": " + e.Err.Error() }

// Unwrap implements the errors.Wrapper interface.
func (e *AuthzError) Unwrap() error { return e.Err }

// NewQueryError returns a new QueryError.
func NewQueryError(op string, err error) *AuthzError {
	return &AuthzError{
		Op:  op,
		Err: err,
	}
}

type storeImpl[W any] struct {
	tx     StoreTxFactory[W]
	client *client.OpenFgaClient
}

// StoreTxFactory is a function that creates a new instance of authz store.
type StoreTxFactory[Tx any] func(*client.OpenFgaClient) (Tx, error)

// NewStore returns a new instance of authz store.
func NewStore[Tx any](client *client.OpenFgaClient, tx StoreTxFactory[Tx]) (Store[Tx], error) {
	return &storeImpl[Tx]{tx, client}, nil
}

// ReadWriteTx starts a read only transaction.
func (s *storeImpl[Tx]) WriteTx(ctx context.Context, fn func(context.Context, Tx) error) error {
	t, err := s.tx(s.client)
	if err != nil {
		return err
	}

	if err := fn(ctx, t); err != nil {
		return err
	}

	return nil
}

type defaultStoreTxImpl struct {
	client *client.OpenFgaClient
}

// DefaultStoreTx is a default authz store transaction.
type DefaultStoreTx interface{}

// NewDefaultStoreTx returns a new instance of default authz store transaction.
func NewDefaultStoreTx() StoreTxFactory[DefaultStoreTx] {
	return func(fga *client.OpenFgaClient) (DefaultStoreTx, error) {
		return &defaultStoreTxImpl{client: fga}, nil
	}
}
