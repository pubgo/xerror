//go:build !cgo

package sqliteutil

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func connectDB(dbPath string) (*sql.DB, error) {
	pragmas := "_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(8000)&_pragma=journal_size_limit(100000000)"
	db, err := sql.Open("sqlite", fmt.Sprintf("%s?%s", dbPath, pragmas))
	if err != nil {
		return nil, err
	}
	return db, nil
}
