package services

import (
	"announce-api/db"
	"announce-api/entities"
	"announce-api/utils"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	client *db.DatabaseClient
}

func (a *Authenticator) CreateUser(user *entities.InputSignUpUser) (int, error) {
	if err := utils.ValidateStruct(user); err != nil {
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("password hashing failed: " + err.Error())
	}

	user.Password = string(hashedPassword)

	if err := validateEmail(user.Email); err != nil {
		return 0, err
	}

	id, err := a.client.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *Authenticator) AuthorizeUser(user *entities.InputSignInUser) (string, error) {
	userFromDb, err := a.client.GetUser(user)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password)); err != nil {
		return "", errors.New("email or password is incorrect")
	}

	tokenStr, err := a.client.GenerateAccessToken(userFromDb)
	if err != nil {
		return "", err
	}

	token, err := a.client.UpdateUserToken(userFromDb, tokenStr)
	if err != nil {
		return "", err
	}

	return token, nil
}

func validateEmail(input string) error {
	matched, err := regexp.Match(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,}$`, []byte(input))
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("email is incorrect")
	}

	return nil
}
