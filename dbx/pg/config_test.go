package pg_test

import (
	"context"
	"testing"

	"github.com/zeiss/pkg/dbx/pg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	config := pg.NewConfig()
	require.NotNil(t, config)
	assert.Equal(t, 5432, config.Port)
	assert.Equal(t, "disable", config.SslMode)
}

func TestFormatDSN(t *testing.T) {
	t.Parallel()

	config := &pg.Config{
		Database: "test_db",
		Host:     "localhost",
		Password: "password",
		Port:     5432,
		SslMode:  "disable",
		User:     "test_user",
	}

	dsn := config.FormatDSN()
	assert.Equal(t, "dbname=test_db user=test_user password=password host=localhost port=5432 sslmode=disable", dsn)
}

func TestContext(t *testing.T) {
	t.Parallel()

	config := pg.Config{
		Database: "test_db",
		Host:     "localhost",
		Password: "password",
		Port:     5432,
		SslMode:  "disable",
		User:     "test_user",
	}

	ctx := config.Context(context.Background())
	assert.Equal(t, config, pg.FromContext(ctx))
}
