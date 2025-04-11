//go:build integration

package postgresql

import (
	"context"
	"fmt"
	"github.com/khostya/pvz/pkg/postgres"
	"github.com/khostya/pvz/pkg/postgres/transactor"
	"os"
	"strings"
)

type DBPool struct {
	pool *postgres.Pool
}

func NewFromEnv() *DBPool {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		panic("TEST_DATABASE_URL isn`t set")
	}

	pool, err := postgres.NewPoolFromURL(context.Background(), url)
	if err != nil {
		panic(err)
	}

	return &DBPool{pool: pool}
}

func (d *DBPool) TruncateTable(ctx context.Context, tableName ...string) {
	if len(tableName) == 0 {
		return
	}
	q := fmt.Sprintf("TRUNCATE %s", strings.Join(tableName, ","))
	if _, err := d.pool.Exec(ctx, q); err != nil {
		panic(err)
	}
}

func (d *DBPool) GetTransactionManager() *transactor.TransactionManager {
	return transactor.NewTransactionManager(d.pool)
}

func (d *DBPool) Close() {
	d.pool.Close()
}
