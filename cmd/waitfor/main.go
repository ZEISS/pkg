package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/zeiss/pkg/cmd/waitfor/schemes"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var build = fmt.Sprintf("%s (%s) (%s)", version, commit, date)

type config struct {
	timeout     time.Duration
	connTimeout time.Duration
	retryTime   time.Duration
}

var cfg = &config{}

var rootCmd = &cobra.Command{
	Use:   "wait-for",
	Short: `wait-for waits for other service to become available.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoot(cmd.Context(), args...)
	},
	Version: build,
}

func init() {
	register(schemes.HTTP(), "http", "https")
	register(schemes.TCP(), "tcp")
	register(schemes.Postgres(), "postgres", "postgresql")

	rootCmd.Flags().DurationVar(&cfg.timeout, "timeout", time.Minute, "Timeout to wait for all checks to complete.")
	rootCmd.Flags().DurationVar(&cfg.connTimeout, "connect-timeout", 5*time.Second, "Timeout to wait for a single check to complete.")
	rootCmd.Flags().DurationVar(&cfg.retryTime, "retry-time", 3*time.Second, "Time to wait between retries.")
}

var waitFuncs = map[string]schemes.WaitFunc{}

func register(fn schemes.WaitFunc, schema ...string) {
	for _, s := range schema {
		waitFuncs[s] = fn
	}
}

func waitFor(ctx context.Context, urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("url parse '%s': %w", urlStr, err)
	}

	t := time.NewTicker(cfg.retryTime)
	defer t.Stop()

	for {
		fn, ok := waitFuncs[u.Scheme]
		if !ok {
			return fmt.Errorf("unsupported schema %q", u.Scheme)
		}

		ct, cancel := context.WithTimeout(ctx, cfg.connTimeout)
		err = fn(ct, urlStr)
		cancel()
		if err == nil {
			return nil
		}

		log.Println("Waiting for", urlStr, err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for %s", urlStr)
		case <-t.C:
		}
	}
}

func runRoot(ctx context.Context, args ...string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.timeout)
	defer cancel()

	for _, urlStr := range args {
		err := waitFor(ctx, urlStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
