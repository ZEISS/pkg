package schemes

import (
	"context"
	"net"
	"net/url"
	"time"
)

// TCP returns a wait function that waits for a TCP connection to be established.
func TCP() WaitFunc {
	return func(ctx context.Context, urlStr string) error {
		u, err := url.Parse(urlStr)
		if err != nil {
			return err
		}

		d, ok := ctx.Deadline()
		if !ok {
			panic("no deadline")
		}

		c, err := net.DialTimeout("tcp", u.Host, time.Until(d))
		if err != nil {
			return err
		}
		defer c.Close()

		return nil
	}
}
