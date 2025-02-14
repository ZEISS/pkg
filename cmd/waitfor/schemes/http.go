package schemes

import (
	"context"
	"fmt"
	"net/http"
)

// HTTP returns a wait function that waits for a HTTP connection to be established.
func HTTP() WaitFunc {
	return func(ctx context.Context, url string) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code %d", resp.StatusCode)
		}

		return nil
	}
}
