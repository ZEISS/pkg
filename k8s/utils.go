package k8s

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IsObjectFound is a helper function to check if an object exists.
func IsObjectFound(ctx context.Context, client client.Client, namespace string, name string, obj client.Object) bool {
	return !errors.IsNotFound(FetchObject(ctx, client, namespace, name, obj))
}

// FetchObject is a helper function to fetch an object.
func FetchObject(ctx context.Context, client client.Client, namespace string, name string, obj client.Object) error {
	return client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, obj)
}

// PhaseExecute executes a function based on a phase.
func PhaseExecute[P comparable](fn func() error, phase P, phases ...P) error {
	for _, p := range phases {
		if phase == p {
			return fn()
		}
	}

	return nil
}
