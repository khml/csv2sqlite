package libc2s

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func ConnectDatabase(dbFilename string) (*sql.DB, error) {
	// connect to SQLite DB
	db, err := sql.Open("sqlite", dbFilename)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
