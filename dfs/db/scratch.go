package db

import "time"

func (d *DB) CreateScratch(name string, year, month, day int, created time.Time) (uint64, error) {
	res, err := d.d.Exec(`
			INSERT INTO scratch (name, created, year, month, day)
			VALUES (?, ?, ?, ?, ?);
		`,
		name, created, year, month, day)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return uint64(id), err
}
