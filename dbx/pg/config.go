package pg

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/zeiss/pkg/utilx"
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

// Config represents configuration for PostgreSQL connection
type Config struct {
	Database string `envconfig:"PG_DB_NAME"`
	Host     string `envconfig:"PG_HOST"`
	Password string `envconfig:"PG_PASSWORD"`
	Port     int    `envconfig:"PG_PORT" default:"5432"`
	SslMode  string `envconfig:"PG_SSL_MODE" default:"disable"`
	User     string `envconfig:"PG_USER"`
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		Port:    5432,
		SslMode: "disable",
	}
}

// FormatDSN formats the given Config into a DSN string which can be passed to the driver.
func (c *Config) FormatDSN() string {
	var params []string

	if utilx.NotEmpty(c.Database) {
		params = append(params, "dbname="+c.escape(c.Database))
	}

	if utilx.NotEmpty(c.User) {
		params = append(params, "user="+c.escape(c.User))
	}

	if utilx.NotEmpty(c.Password) {
		params = append(params, "password="+c.escape(c.Password))
	}

	if utilx.NotEmpty(c.Host) {
		params = append(params, "host="+c.escape(c.Host))
	}

	if utilx.NotEmpty(c.Port) {
		params = append(params, "port="+strconv.Itoa(c.Port))
	}

	if utilx.NotEmpty(c.SslMode) {
		params = append(params, "sslmode="+c.escape(c.SslMode))
	}

	return strings.Join(params, " ")
}

func (c *Config) escape(str string) string {
	if !strings.Contains(str, " ") {
		return str
	}

	str = strings.ReplaceAll(str, "'", "\\'")
	return fmt.Sprintf("'%s'", str)
}
