//go:build integration

package postgres

import (
	"context"
	"github.com/khostya/pvz/internal/repo/postgres/schema"
	"github.com/khostya/pvz/tests/postgres/postgresql"
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

	truncate()
	code := m.Run()

	db.Close()

	os.Exit(code)
}

func truncate() {
	db.TruncateTable(context.Background(), tables...)
}
