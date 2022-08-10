package sqliteutil

import "database/sql"

func ConnectDB(dbPath string) (*sql.DB, error) {
	return connectDB(dbPath)
}
