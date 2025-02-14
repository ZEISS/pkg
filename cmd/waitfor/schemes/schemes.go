package schemes

import "context"

type WaitFunc func(context.Context, string) error
