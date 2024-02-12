package types

type ID string

type User struct {
	ID            ID
	Login         string
	Email         string
	Password      string
	IsAdmin       bool
	HiddenAuthors []ID
}
