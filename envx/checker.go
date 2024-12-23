package envx

import "context"

// Check is the function signature for the checker.
type Check func(ctx context.Context) error

// Checker is the interface for the checker.
type Checker interface {
	// Check runs the checker.
	Check(ctx context.Context, check ...Check) error
}

type checker struct{}

// NewChecker returns a new checker.
func NewChecker() *checker {
	return &checker{}
}

// Check runs the checker.
func (c *checker) Check(ctx context.Context, check ...Check) error {
	for _, c := range check {
		if err := c(ctx); err != nil {
			return err
		}
	}

	return nil
}
