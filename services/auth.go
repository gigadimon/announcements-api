package services

import (
	"announce-api/db"
	"announce-api/entities"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	client *db.DatabaseClient
}

func (a *Authenticator) CreateUser(user *entities.InputUser) (bool, error) {
	if user.Email == "" {
		return false, errors.New("email required")
	}

	if user.Login == "" {
		return false, errors.New("login required")
	}

	if user.Password == "" {
		return false, errors.New("password required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, errors.New("password hashing failed: " + err.Error())
	}

	user.Password = string(hashedPassword)

	if err := validateEmail(user.Email); err != nil {
		return false, err
	}

	if err := a.client.CreateUser(user); err != nil {
		return false, err
	}

	return false, nil
}

func (a *Authenticator) AuthorizeUser(user *entities.InputUser) {}

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
