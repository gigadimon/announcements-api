package entities

import "time"

type ID int

type User struct {
	ID       ID     `db:"id"`
	Login    string `db:"login"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type InputUser struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
