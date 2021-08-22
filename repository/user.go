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
	stmt, err := ur.db.Prepare("select * from userpet where email = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(email).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepositorySql) FindById(id string) (*model.User, error) {
	stmt, err := ur.db.Prepare("select * from userpet where id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepositorySql) Update(user *model.User) error {
	stmt, err := ur.db.Prepare("update userpet set name = $1, password_hash = $2 where id = $3")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.PasswordHash, user.Id)
	if err != nil {
		return err
	}

	return nil
}
