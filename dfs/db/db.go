package db

import (
	"os"

	"database/sql"

	// Import the sqlite3 driver for package sql
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	d *sql.DB
}

func Open(dbPath string) (*DB, error) {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		err = migrateDb(dbPath)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	d, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &DB{
		d: d,
	}, nil
}

func (d *DB) Close() error {
	return d.d.Close()
}
