package fga

import (
	"context"
	"strings"

	"github.com/openfga/go-sdk/client"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/conv"
)

// User is the user that is making the request.
type User string

// Object is the object that is being accessed.
type Object string

// Relation is the relation between the user and the object.
type Relation string

// NoopUser is a user that represents no user.
const NoopUser User = ""

// NoopRelation is a relation that represents no relation.
const NoopRelation Relation = ""

// NoopObject is an object that represents no object.
const NoopObject Object = ""

// Stringer create a string an adds it to the representation.
type Stringer func() string

// Store is an interface that provides methods for transactional operations on the authz database.
type Store[Tx any] interface {
	// Allowed checks if the user is allowed to perform the operation on the object.
	Allowed(context.Context, User, Object, Relation) (bool, error)
	// WriteTx starts a read write transaction.
	WriteTx(context.Context, func(context.Context, Tx) error) error
}

// StoreTx is an interface that provides methods for transactional operations on the authz database.
type StoreTx interface {
	// WriteTuple writes a tuple to the authz database.
	WriteTuple(context.Context, User, Object, Relation) error
	// DeleteTuple deletes a tuple from the authz database.
	DeleteTuple(context.Context, User, Object, Relation) error
}

// NoopStore is a store that does nothing.
type NoopStore struct{}

// Allowed checks if the user is allowed to perform the operation on the object.
func (n *NoopStore) Allowed(context.Context, User, Object, Relation) (bool, error) {
	return true, nil
}

// WriteTx starts a read write transaction.
func (n *NoopStore) WriteTx(context.Context, func(context.Context, StoreTx) error) error {
	return nil
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
type StoreTxFactory[Tx any] func(*client.OpenFgaClient, StoreTx) (Tx, error)

// NewStore returns a new instance of authz store.
func NewStore[Tx any](client *client.OpenFgaClient, tx StoreTxFactory[Tx]) (Store[Tx], error) {
	return &storeImpl[Tx]{tx, client}, nil
}

// Allowed checks if the user is allowed to perform the operation on the object.
func (t *storeImpl[W]) Allowed(ctx context.Context, user User, object Object, relation Relation) (bool, error) {
	opts := client.ClientCheckOptions{}

	body := client.ClientCheckRequest{
		User:     conv.String(user),
		Relation: conv.String(relation),
		Object:   conv.String(object),
	}

	data, err := t.client.Check(ctx).Options(opts).Body(body).Execute()
	if err != nil {
		return false, err
	}

	ok := cast.Value(data.Allowed)
	if ok {
		return true, nil
	}

	return false, nil
}

// ReadWriteTx starts a read only transaction.
func (s *storeImpl[Tx]) WriteTx(ctx context.Context, fn func(context.Context, Tx) error) error {
	t, err := s.tx(s.client, s)
	if err != nil {
		return err
	}

	if err := fn(ctx, t); err != nil {
		return err
	}

	return nil
}

// WriteTuple writes a tuple to the authz database.
func (s *storeImpl[Tx]) WriteTuple(ctx context.Context, user User, object Object, relation Relation) error {
	body := client.ClientWriteRequest{
		Writes: []client.ClientTupleKey{
			{
				User:     conv.String(user),
				Relation: conv.String(relation),
				Object:   conv.String(object),
			},
		},
	}

	_, err := s.client.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

// DeleteTuple deletes a tuple from the authz database.
func (s *storeImpl[Tx]) DeleteTuple(ctx context.Context, user User, object Object, relation Relation) error {
	body := client.ClientWriteRequest{
		Deletes: []client.ClientTupleKeyWithoutCondition{
			{
				User:     conv.String(user),
				Relation: conv.String(relation),
				Object:   conv.String(object),
			},
		},
	}

	_, err := s.client.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

// DefaultSeparator is the default separator for entities.
const DefaultSeparator = "/"

// DefaultNamespaceSeparator is the default separator for namespaces.
const DefaultNamespaceSeparator = ":"

// EntitiesString is a type that represents a list of entities.
func EntityString[E Entities](e E) string {
	return conv.String(e)
}

// Entities is a type that represents a list of entities.
type Entities interface {
	User | Relation | Object
}

// NewEntity returns a new User.
func NewEntity[E Entities](s ...Stringer) E {
	u := ""

	for _, v := range s {
		u += v()
	}

	return E(u)
}

// NewUser returns a new User.
func NewUser(s ...Stringer) User {
	return NewEntity[User](s...)
}

// NewRelation returns a new Relation.
func NewRelation(s ...Stringer) Relation {
	return NewEntity[Relation](s...)
}

// NewObject returns a new Object.
func NewObject(s ...Stringer) Object {
	return NewEntity[Object](s...)
}

// Namespace adds a namespace to the entity.
func Namespace(namespace string, sep ...string) Stringer {
	return func() string {
		s := DefaultNamespaceSeparator

		if len(sep) > 0 {
			s = sep[0]
		}

		if namespace == "" {
			return ""
		}

		return namespace + s
	}
}

// Join joins the entities with the separator.
func Join(sep string, entities ...string) Stringer {
	return func() string {
		return strings.Join(entities, sep)
	}
}

// String returns the string representation of the entity.
func String(s string) Stringer {
	return func() string {
		return s
	}
}
