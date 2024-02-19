package services

import (
	"announce-api/db"
	"announce-api/entities"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	client *db.DatabaseClient
}

func (a *Authenticator) CreateUser(user *entities.InputSignUpUser) (int, error) {
	if err := user.Validate(); err != nil {
		return 0, errors.New("passed user is incorrect: " + err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("password hashing failed: " + err.Error())
	}

	user.Password = string(hashedPassword)

	id, err := a.client.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *Authenticator) AuthorizeUser(user *entities.InputSignInUser) (string, error) {
	if err := user.Validate(); err != nil {
		return "", errors.New("passed user is incorrect: " + err.Error())
	}

	userFromDb, err := a.client.GetUser(user)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password)); err != nil {
		return "", errors.New("login or password is incorrect")
	}

	tokenStr, err := generateAccessToken(userFromDb)
	if err != nil {
		return "", err
	}

	token, err := a.client.UpdateUserToken(userFromDb, tokenStr)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *Authenticator) IsTokenExists(token string) (int, error) {
	id, err := a.client.IsTokenExists(token)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func generateAccessToken(user *entities.User) (string, error) {
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
