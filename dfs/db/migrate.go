package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func migrateDb(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE tag (
			tag_id INTEGER PRIMARY KEY,
			name TEXT,
			UNIQUE(name)
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE year (
			year INTEGER,
			PRIMARY KEY(year)
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE month (
			year INTEGER,
			month INTEGER,
			PRIMARY KEY(year, month),
			FOREIGN KEY(year) REFERENCES year(year)
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE day (
			year INTEGER,
			month INTEGER,
			day INTEGER,
			PRIMARY KEY(year, month, day),
			FOREIGN KEY(year, month) REFERENCES month(year, month)
		);
	`)
	if err != nil {
		return err
	}

	// Docs
	/*
		id
		year
		month
		day
		uuid
		filename
		checksum
	*/

	// Doc Tags
	/*
		tag_id
		doc_id
	*/

	return nil
}
