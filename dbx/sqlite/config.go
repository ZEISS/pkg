package sqlite

import (
	"context"
	"path/filepath"

	"github.com/zeiss/pkg/filex"
)

type contextKey int

const (
	configKey contextKey = iota
)

// Context returns a new Context that carries the provided Config.
func (cfg *Config) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

// FromContext will return the Config carried in the provided Context.
//
// It panics if config is not available on the current context.
func FromContext(ctx context.Context) *Config {
	return ctx.Value(configKey).(*Config)
}

// Config represents configuration for PostgreSQL connection.
type Config struct {
	// Path is the path to the SQLite database file.
	Path string `envconfig:"DBX_SQLITE_PATH"`
}

// NewConfig returns a new Config instance.
func NewConfig(path string) *Config {
	return &Config{Path: path}
}

// Dir returns the directory of the SQLite database file.
func (c *Config) Dir() (string, error) {
	return filepath.Abs(filepath.Dir(c.Path))
}

// Mkdir creates the directory for the SQLite database file.
func (c *Config) Mkdir() error {
	dir, err := c.Dir()
	if err != nil {
		return err
	}

	return filex.MkdirAll(dir, 0o755)
}
