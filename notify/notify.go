package notify

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

// ErrSendNotification is returned when a notification could not be sent.
var ErrSendNotification = errors.New("notify: could not send notification")

// Config is the configuration for a notifier.
type Config struct {
	DeviceTokens []string
}

// Notifier is an interface for sending notifications.
type Notifier interface {
	Notify(ctx context.Context, title, message string, config ...Config) error
}

var _ Notifier = (*Notify)(nil)

// Notify sends a notification with the given title and message.
type Notify struct {
	notifiers []Notifier
}

// Opt is a functional option for configuring Notify.
type Opt func(*Notify)

// WithNotifiers sets the notifiers to use.
func WithNotifiers(notifiers ...Notifier) Opt {
	return func(n *Notify) {
		n.notifiers = notifiers
	}
}

// New creates a new Notify.
func New(opts ...Opt) *Notify {
	n := &Notify{
		notifiers: []Notifier{},
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

func (n *Notify) send(ctx context.Context, subject, message string, config ...Config) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, service := range n.notifiers {
		if service == nil {
			continue
		}

		g.Go(func() error {
			return service.Notify(ctx, subject, message, config...)
		})
	}

	err := g.Wait()
	if err != nil {
		err = errors.Join(ErrSendNotification, err)
	}

	return err
}

// Notify sends a notification with the given title and message.
func (n *Notify) Notify(ctx context.Context, title, message string, config ...Config) error {
	return n.send(ctx, title, message, config...)
}
