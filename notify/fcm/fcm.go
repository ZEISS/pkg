package fcm

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeiss/pkg/notify"

	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/fcm/v1"
)

var (
	// ErrNoConfig is returned when no configuration is provided.
	ErrNoConfig = errors.New("fcm: no configuration provided")
	// ErrNoDeviceTokens is returned when no device tokens are provided.
	ErrNoDeviceTokens = errors.New("fcm: no device tokens provided")
)

// Compile-time check that Service satisfies the Notifier interface.
var _ notify.Notifier = (*FCM)(nil)

type fcmClient interface {
	Send(ctx context.Context, message ...*messaging.Message) (*messaging.BatchResponse, error)
	SendMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
}

// FCM is a notifier for Firebase Cloud Messaging.
type FCM struct {
	client fcmClient
}

// Opt is a functional option for configuring FCM.
type Opt func(*FCM)

// New creates a new FCM.
func New(client fcmClient, opts ...Opt) *FCM {
	f := &FCM{
		client: client,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

// Notify sends a notification with the given title and message.
func (f *FCM) Notify(ctx context.Context, title, message string, config ...notify.Config) error {
	if len(config) == 0 {
		return ErrNoConfig
	}

	fcm.NewService(ctx)
	cfg := config[0]
	if len(cfg.DeviceTokens) == 0 {
		return ErrNoDeviceTokens
	}

	msg := &messaging.MulticastMessage{
		Tokens: cfg.DeviceTokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
	}

	_, err := f.client.SendMulticast(ctx, msg)
	if err != nil {
		return fmt.Errorf("fcm: send multicast message to FCM devices: %w", err)
	}

	return nil
}
