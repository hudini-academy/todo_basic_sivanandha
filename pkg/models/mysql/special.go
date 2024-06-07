package mysql

import (
	"TODO/pkg/models"
	"database/sql"
)

type SpecialModel struct {
	DB *sql.DB
}

func (m *SpecialModel) Insert(name, description,expires string) (int, error) {
	stmt := `INSERT INTO specials (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	//  Use the Exec() method on the embedded connection pool to execute the statement
	result, err := m.DB.Exec(stmt, name, description,expires)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the todos table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type before returning.
	return int(id), nil
}
func (m *SpecialModel) Latest() ([]*models.Special, error) {
	stmt := `SELECT id, title,content, created, expires FROM specials
	WHERE expires > UTC_TIMESTAMP()`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	specialss := []*models.Special{}
	for rows.Next(){
		s := &models.Special{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		specialss= append(specialss, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return specialss, nil
}
func (m *SpecialModel) Delete(id int) error{
	// SQL statement for deleting task
	stmt:= `DELETE FROM specials WHERE id = ?;`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
  		panic(err)
	}
	return nil
}