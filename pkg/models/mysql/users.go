package mysql

import (
	"TODO/pkg/models"
	"database/sql"
	// "strings"
	// "github.com/go-sql-driver/mysql"
	// "golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	stmnt := `INSERT INTO users(name,email, hashed_password, created)
	VALUES (?,?,?,UTC_TIMESTAMP())`
	_, err := m.DB.Exec(stmnt, name, email, password)
	if err != nil {
		return err
	}
	return err
}
func (m *UserModel) Authenticate(email, password string) (bool, error) {
	stmt := `SELECT id from users WHERE email=? AND hashed_password=?`
	rows, err := m.DB.Query(stmt, email, password)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	return rows.Next(), nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
