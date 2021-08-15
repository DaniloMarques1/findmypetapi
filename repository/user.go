package repository

import (
	"database/sql"
	"github.com/danilomarques1/findmypetapi/model"
)

type UserRepositorySql struct {
	db *sql.DB
}

func NewUserRepositorySql(db *sql.DB) *UserRepositorySql {
	return &UserRepositorySql{
		db: db,
	}
}

func (ur *UserRepositorySql) Save(user *model.User) error {
	stmt, err := ur.db.Prepare("INSERT INTO userpet(id, name, email, password_hash) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id, user.Name, user.Email, user.PasswordHash); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepositorySql) FindByEmail(email string) (*model.User, error) {
	return nil, nil
}

func (ur *UserRepositorySql) FindById(id string) (*model.User, error) {
	return nil, nil
}

func (ur *UserRepositorySql) Update(user *model.User) error {
	return nil
}
