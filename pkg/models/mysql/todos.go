package mysql

import (
	"TODO/pkg/models"
	"database/sql"
)

// define the TodosModel type it wrap the db connection pool
type TodoModel struct{
	DB *sql.DB
}
// This will insert a new todos into the database.
func (m *TodoModel) Insert(name,description, expires string) (int, error) {
	//SQL statement we want to execute
	stmt := `INSERT INTO todos (name, description, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	//  Use the Exec() method on the embedded connection pool to execute the statement
	result , err := m.DB.Exec(stmt,name, description, expires)
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
// created a function for delete task
func (m *TodoModel) Delete(title string) error{
	// SQL statement for deleting task
	var exists bool
	stmt:= `DELETE FROM todos WHERE name = ?`
	stmt2 := `DELETE FROM specials WHERE title = ?`
	errexis := m.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM specials WHERE title = ?)",title).Scan(&exists)
	if errexis!= nil{
		panic(errexis)

	}

	_, err := m.DB.Exec(stmt, title)
	if err != nil {
  		panic(err)
	}
	if exists{
	_, err1 := m.DB.Exec(stmt2, title)
	if err1 != nil {
  		panic(err1)
	}
}
	return nil
}
// This function will return a specific todos based on its id.
func (m *TodoModel) Get(id int) (*models.Todos, error) {
	// SQL statement
	stmt := `SELECT * FROM todos WHERE id = ?`
	// Use the QueryRow() method on the connection pool to execute stmt
	row := m.DB.QueryRow(stmt, id)
	// Initialize a pointer to a new zeroed todos struct.
	s := &models.Todos{}
	// Use row.Scan() to copy the values from each field in sql.Row to the 
	// corresponding field in the Todos struct
	err := row.Scan(&s.ID, &s.Name, &s.Created, &s.Expires, &s.Description)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}
// This will return the 10 most recently created todos.
func (m *TodoModel) GetAllTask() ([]*models.Todos, error) {
	stmt := `SELECT id, name, description, created, expires FROM todos
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todoos := []*models.Todos{}
	for rows.Next(){
		s := &models.Todos{}
		err = rows.Scan(&s.ID, &s.Name, &s.Description, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of todoos
		todoos= append(todoos, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todoos, nil
}
// this function will update the task 
func (m *TodoModel) Update(id int, name string) (bool, error){
	stmnt := `UPDATE todos SET name=? WHERE id=?`
	_, err:= m.DB.Exec(stmnt,name, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

	