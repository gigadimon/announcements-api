package entities

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
)

type PhotosForDB interface {
	driver.Valuer
	sql.Scanner
}

type InputSignUpUser struct {
	Login    string `json:"login" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type InputSignInUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type InputAnnouncement struct {
	AuthorPhone string `form:"author_phone" json:"author_phone"`
	Title       string `form:"title" json:"title" validate:"required"`
	Description string `form:"description" json:"description" validate:"required"`
}

type AnnouncementToDB struct {
	ID          int         `db:"id" json:"id"`
	AuthorID    int         `db:"author_id" json:"author_id"`
	AuthorLogin string      `db:"author_login" json:"author_login"`
	AuthorEmail string      `db:"author_email" json:"author_email"`
	AuthorPhone string      `db:"author_phone" json:"author_phone"`
	Title       string      `db:"title" json:"title" validate:"required"`
	Photos      PhotosForDB `db:"photos" json:"photos"`
	Description string      `db:"description" json:"description" validate:"required"`
	CreatedAt   time.Time   `db:"created_at" json:"created_at"`
}

type AnnouncementFromDB struct {
	ID          int            `db:"id" json:"id"`
	AuthorID    int            `db:"author_id" json:"author_id"`
	AuthorLogin string         `db:"author_login" json:"author_login"`
	AuthorEmail string         `db:"author_email" json:"author_email"`
	AuthorPhone string         `db:"author_phone" json:"author_phone"`
	Title       string         `db:"title" json:"title" validate:"required"`
	Photos      pq.StringArray `db:"photos" json:"photos"`
	Description string         `db:"description" json:"description" validate:"required"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	IsHidden    bool           `db:"is_hidden" json:"is_hidden"`
}

type AuthorInfo struct {
	ID    int
	Login string
	Email string
}
