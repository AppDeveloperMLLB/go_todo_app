package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	port := 5432
	if _, defined := os.LookupEnv("CI"); defined {
		port = 5432
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=127.0.0.1 user=todo password=todo dbname=todo port=%d sslmode=disable",
		port,
	))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return sqlx.NewDb(db, "postgres")
}
