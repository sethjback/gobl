package sqlite

import (
	"database/sql"
	"errors"

	"github.com/sethjback/gobl/model"
)

func (d *SQLite) GetUser(email string) (*model.User, error) {
	var password string
	var lastlogin interface{}
	err := d.Connection.QueryRow("SELECT email, password, lastlogin FROM "+usersTable+" where email=?", email).Scan(&email, &password, &lastlogin)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Could not find that user")
	case err != nil:
		return nil, err
	default:
		return &model.User{Email: email, Password: password}, nil
	}
}

func (d *SQLite) UserList() ([]model.User, error) {
	rows, err := d.Connection.Query("SELECT email, password FROM " + usersTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User

	for rows.Next() {
		var email, password string
		err = rows.Scan(&email, &password)
		if err != nil {
			return nil, err
		}
		users = append(users, model.User{Email: email, Password: password})
	}
	return users, nil
}

func (d *SQLite) SaveUser(u model.User) error {
	_, err := d.GetUser(u.Email)
	if err != nil {
		if err.Error() != "Could not find that user" {
			return err
		}
		return insertUser(d, u)
	}
	return updateUser(d, u)
}

func updateUser(d *SQLite, u model.User) error {
	sql := "UPDATE " + usersTable + " set email=?, password=? WHERE email=?"
	_, err := d.Connection.Exec(sql, u.Email, u.Password, u.Email)
	return err
}

func insertUser(d *SQLite, u model.User) error {
	sql := "INSERT INTO " + usersTable + " (email, password) values (?, ?)"

	_, err := d.Connection.Exec(sql, u.Email, u.Password)

	return err
}

func (d *SQLite) DeleteUser(email string) error {
	_, err := d.Connection.Exec("DELETE FROM "+usersTable+" WHERE email=?", email)

	return err
}
