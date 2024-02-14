package db

import (
	"announce-api/entities"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (a *Authenticator) CreateUser(user *entities.InputSignUpUser) (int, error) {
	var id int
	query := `INSERT INTO users (email, login, password) VALUES ($1, $2, $3) RETURNING id`

	row := a.db.QueryRow(query, user.Email, user.Login, user.Password)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (a *Authenticator) GetUser(user *entities.InputSignInUser) (*entities.User, error) {
	userFromDb := new(entities.User)
	selectQuery := "SELECT id, email, login, password FROM users WHERE email=$1"

	if err := a.db.Get(userFromDb, selectQuery, user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(entities.User), errors.New("user with passed email doesn't exists")
		}
		return new(entities.User), errors.New("error while getting info from database: " + err.Error())
	}

	return userFromDb, nil
}

func (a *Authenticator) UpdateUserToken(user *entities.User, token string) (string, error) {
	updateQuery := `UPDATE users SET token=$1 WHERE email=$2 RETURNING token`

	row := a.db.QueryRow(updateQuery, token, user.Email)

	var writtenToken string
	if err := row.Scan(&writtenToken); err != nil {
		return "", err
	}

	return writtenToken, nil
}

func (a *Authenticator) GenerateAccessToken(user *entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"login":   user.Login,
		"expires": time.Now().Add(time.Hour * 10).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
