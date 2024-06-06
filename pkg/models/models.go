package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
// this will be use later when user give incorrect details
var ErrInvalidCredentials = errors.New("models: invalid credentials")
// this will be use when user give the email that alrdy used
var ErrDuplicateEmail = errors.New("models: duplicate email")

type Todos struct{
	ID int
	Name string
	Description string
	Created time.Time
	Expires time.Time
}
type User struct {
	ID int
	UserName string
	Email string
	HashedPassword []byte
	Created time.Time
}
	