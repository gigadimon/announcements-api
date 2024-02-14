package entities

import "time"

type ID int

type InputSignUpUser struct {
	Login    string `json:"login" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type InputSignInUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID       ID     `db:"id"`
	Login    string `db:"login"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type Announcement struct {
	ID          ID     `db:"id"`
	AuthorID    ID     `db:"author_id"`
	AuthorLogin string `db:"author_login"`
	AuthorEmail string `db:"author_email"`
	AuthorPhone string `db:"author_phone"`
	Title       string `db:"title"`
	// Photos      []string  `db:"author_id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}
