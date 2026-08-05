package storex

import (
	"context"
	"time"
)

// StorageWithConn interface for communicating with different database/key-value
// providers that require a connection.
type StorageWithConn[T any] interface {
	// Conn returns a connection to the storage.
	Conn() T
}

// Store interface for communicating with different database/key-value
// providers.
type Store interface {
	// GetWithContext gets the value for the given key with a context.
	// `nil, nil` is returned when the key does not exist
	GetWithContext(ctx context.Context, key string) ([]byte, error)
	// Get gets the value for the given key.
	// `nil, nil` is returned when the key does not exist
	Get(key string) ([]byte, error)
	// SetWithContext stores the given value for the given key
	// with an expiration value, 0 means no expiration.
	SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error
	// Set stores the given value for the given key along
	// with an expiration value, 0 means no expiration.
	// Empty key or value will be ignored without an error.
	Set(key string, val []byte, exp time.Duration) error
	// DeleteWithContext deletes the value for the given key with a context.
	// It returns no error if the storage does not contain the key,
	DeleteWithContext(ctx context.Context, key string) error
	// Delete deletes the value for the given key.
	// It returns no error if the storage does not contain the key,
	Delete(key string) error
	// ResetWithContext resets the storage and deletes all keys with a context.
	ResetWithContext(ctx context.Context) error
	// Reset resets the storage and delete all keys.
	Reset() error
	// Close closes the storage and will stop any running garbage
	// collectors and open connections.
	Close() error
}
