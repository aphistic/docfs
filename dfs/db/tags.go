package db

type Tag struct {
	ID   uint64
	Name string
}

func (d *DB) GetTags() ([]*Tag, error) {
	tags := make([]*Tag, 0)

	res, err := d.d.Query("SELECT tag_id, name FROM tag")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var tagID uint64
		var tagName string
		err = res.Scan(&tagID, &tagName)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &Tag{
			ID:   tagID,
			Name: tagName,
		})
	}

	return tags, nil
}

func (d *DB) GetTag(tag string) (*Tag, error) {
	res, err := d.d.Query("SELECT tag_id, name FROM tag WHERE name == ?", tag)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if !res.Next() {
		return nil, nil
	}

	var tagID uint64
	var tagName string
	err = res.Scan(&tagID, &tagName)
	if err != nil {
		return nil, err
	}

	return &Tag{
		ID:   tagID,
		Name: tagName,
	}, nil
}

func (d *DB) AddTag(tag string) (uint64, error) {
	var tagID uint64
	res, err := d.d.Query("SELECT tag_id FROM tag WHERE name == ?", tag)
	if err != nil {
		return 0, err
	}
	defer res.Close()
	if !res.Next() {
		res, err := d.d.Exec("INSERT INTO tag (name) VALUES (?)", tag)
		if err != nil {
			return 0, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		return uint64(id), nil
	}

	err = res.Scan(&tagID)
	if err != nil {
		return 0, err
	}

	return uint64(tagID), nil
}

func (d *DB) RemoveTag(tag string) error {
	_, err := d.d.Exec("DELETE FROM tag WHERE name == ?", tag)
	if err != nil {
		return err
	}
	return nil
}
