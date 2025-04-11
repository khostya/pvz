//go:build integration

package postgres

import (
	"context"
	"github.com/khostya/pvz/internal/repo/postgres/schema"
	"github.com/khostya/pvz/test/postgres/postgresql"
	"os"
	"testing"
)

var (
	db *postgresql.DBPool
)

var tables []string = []string{
	schema.Product{}.TableName(),
	schema.Reception{}.TableName(),
	schema.PVZ{}.TableName(),
	schema.User{}.TableName(),
}

func TestMain(m *testing.M) {
	db = postgresql.NewFromEnv()

	code := m.Run()

	db.TruncateTable(context.Background(), tables...)
	db.Close()

	os.Exit(code)
}
