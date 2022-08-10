//go:build cgo

package sqliteutil

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func connectDB(dbPath string) (*sql.DB, error) {
	pragmas := "_foreign_keys=1&_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=8000"

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?%s", dbPath, pragmas))
	if err != nil {
		return nil, err
	}

	// additional pragmas not supported through the dsn string
	_, err = db.Exec("pragma journal_size_limit = 100000000;")
	if err != nil {
		return nil, err
	}

	return db, err
}
