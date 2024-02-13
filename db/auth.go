package db

import (
	"announce-api/entities"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Authenticator struct {
	db *sqlx.DB
}

type NewUser struct {
	ID    int
	Email string
	Login string
}

func (a *Authenticator) CreateUser(user *entities.InputUser) error {
	rows, err := a.db.NamedQuery(`INSERT INTO users (email, login, password) VALUES (:email, :login, :password) RETURNING id, email, login`, user)
	if err != nil {
		return err
	}

	nu := new(NewUser)

	defer rows.Close()
	for rows.Next() {
		err = rows.StructScan(&nu)
		if err != nil {
			return errors.New("Row scan failed: " + err.Error())
		}
	}

	return nil
}

func (a *Authenticator) AuthorizeUser() {

}

func (a *Authenticator) generateUserAccessToken() {
}
