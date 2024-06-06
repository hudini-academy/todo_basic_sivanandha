package mysql

import (
	"TODO/pkg/models"
	"database/sql"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashed_password, errr := bcrypt.GenerateFromPassword([]byte(password), 12)
	if errr != nil {
		return errr
	}
	stmnt := `INSERT INTO users(name,email, hashed_password, created)
	VALUES (?,?,?,UTC_TIMESTAMP())`
	_, err := m.DB.Exec(stmnt,name,email,string(hashed_password))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "u"){ 
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `SELECT id from users WHERE email=? AND password=?`
	_,err := m.DB.Exec(stmt,email,password)
	if err != nil {
		return 0, nil
	}
	return 0, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}