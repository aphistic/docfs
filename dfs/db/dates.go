package db

// Years

func (d *DB) GetYears() ([]uint64, error) {
	res, err := d.d.Query("SELECT year FROM year;")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	years := make([]uint64, 0)
	for res.Next() {
		var year uint64
		err = res.Scan(&year)
		if err != nil {
			return nil, err
		}
		years = append(years, year)
	}

	return years, nil
}

func (d *DB) GetYear(year uint64) (uint64, error) {
	res, err := d.d.Query("SELECT year FROM year WHERE year == ?", year)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	if res.Next() {
		return year, nil
	}

	return 0, ErrNotExists
}

func (d *DB) AddYear(year uint64) error {
	res, err := d.d.Query("SELECT year FROM year WHERE year == ?", year)
	if err != nil {
		return err
	}
	defer res.Close()

	if !res.Next() {
		_, err := d.d.Exec("INSERT INTO year (year) VALUES (?);", year)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) RemoveYear(year uint64) error {
	_, err := d.d.Exec("DELETE FROM year WHERE year == ?", year)
	if err != nil {
		return err
	}
	_, err = d.d.Exec("DELETE FROM month WHERE year == ?", year)
	if err != nil {
		return err
	}
	_, err = d.d.Exec("DELETE FROM day WHERE year == ?", year)
	if err != nil {
		return err
	}

	return nil
}

// Months

func (d *DB) GetMonths(year uint64) ([]uint64, error) {
	res, err := d.d.Query("SELECT month FROM month WHERE year == ?", year)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	months := make([]uint64, 0)
	for res.Next() {
		var month uint64
		err = res.Scan(&month)
		if err != nil {
			return nil, err
		}
		months = append(months, month)
	}

	return months, nil
}

func (d *DB) GetMonth(year uint64, month uint64) (uint64, error) {
	res, err := d.d.Query("SELECT year, month FROM month WHERE year == ? AND month == ?", year, month)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	if res.Next() {
		return month, nil
	}

	return 0, ErrNotExists
}

func (d *DB) AddMonth(year uint64, month uint64) error {
	err := d.AddYear(year)
	if err != nil {
		return err
	}

	res, err := d.d.Query("SELECT month, year FROM month WHERE year == ? AND month == ?", year, month)
	if err != nil {
		return err
	}
	defer res.Close()

	if !res.Next() {
		_, err := d.d.Exec("INSERT INTO month (year, month) VALUES (?, ?);", year, month)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) RemoveMonth(year uint64, month uint64) error {
	_, err := d.d.Exec("DELETE FROM month WHERE year == ? AND month == ?", year, month)
	if err != nil {
		return err
	}

	return nil
}

// Days

func (d *DB) GetDays(year uint64, month uint64) ([]uint64, error) {
	res, err := d.d.Query("SELECT day FROM day WHERE year == ? AND month == ?", year, month)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	days := make([]uint64, 0)
	for res.Next() {
		var day uint64
		err = res.Scan(&day)
		if err != nil {
			return nil, err
		}
		days = append(days, day)
	}

	return days, nil
}

func (d *DB) GetDay(year uint64, month uint64, day uint64) (uint64, error) {
	res, err := d.d.Query("SELECT year, month, day FROM day WHERE year == ? AND month == ? AND day == ?", year, month, day)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	if res.Next() {
		return day, nil
	}

	return 0, ErrNotExists
}

func (d *DB) AddDay(year uint64, month uint64, day uint64) error {
	err := d.AddMonth(year, month)
	if err != nil {
		return err
	}

	res, err := d.d.Query("SELECT year, month, day FROM day WHERE year == ? AND month == ? AND day == ?", year, month, day)
	if err != nil {
		return err
	}
	defer res.Close()

	if !res.Next() {
		_, err := d.d.Exec("INSERT INTO day (year, month, day) VALUES (?, ?, ?);", year, month, day)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) RemoveDay(year uint64, month uint64, day uint64) error {
	_, err := d.d.Exec("DELETE FROM day WHERE year == ? AND month == ? AND day == ?", year, month, day)
	if err != nil {
		return err
	}

	return nil
}
