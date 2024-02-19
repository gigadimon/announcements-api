package db

import (
	"announce-api/entities"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Authenticator struct {
	db *sqlx.DB
}

type NewUser struct {
	ID    int
	Email string
	Login string
}

func (a *Authenticator) CreateUser(user *entities.InputSignUpUser) (int, error) {
	var id int
	query := `INSERT INTO users (email, login, password) VALUES ($1, $2, $3) RETURNING id`
	rows, err := a.db.Query(query, user.Email, user.Login, user.Password)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			return 0, errors.New(pqError.Detail)
		}
		return 0, err
	}

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (a *Authenticator) GetUser(user *entities.InputSignInUser) (*entities.User, error) {
	userFromDb := new(entities.User)
	selectQuery := "SELECT id, email, login, password FROM users WHERE login=$1"

	if err := a.db.Get(userFromDb, selectQuery, user.Login); err != nil {
		return new(entities.User), errors.New("error while getting info from database: " + err.Error())
	}

	return userFromDb, nil
}

func (a *Authenticator) UpdateUserToken(user *entities.User, token string) (string, error) {
	updateQuery := `UPDATE users SET token=$1 WHERE login=$2 RETURNING token`

	row := a.db.QueryRow(updateQuery, token, user.Login)

	var writtenToken string
	if err := row.Scan(&writtenToken); err != nil {
		return "", err
	}

	return writtenToken, nil
}

func (a *Authenticator) IsTokenExists(token string) (int, error) {
	var id int
	selectQuery := "SELECT id FROM users WHERE token=$1"
	if err := a.db.Get(&id, selectQuery, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("token invalid")
		}
		return 0, errors.New("error while getting info from database: " + err.Error())
	}

	return id, nil
}
