package db

import (
	"announce-api/entities"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AnnouncementManager struct {
	db *sqlx.DB
}

func (m *AnnouncementManager) GetList() {}

func (m *AnnouncementManager) GetOneById(id string) (*entities.AnnouncementFromDB, error) {
	announcement := new(entities.AnnouncementFromDB)
	query := `SELECT * FROM announcements WHERE id=$1`

	row := m.db.QueryRowx(query, id)
	if err := row.StructScan(announcement); err != nil {
		return new(entities.AnnouncementFromDB), err
	}
	return announcement, nil
}

func (m *AnnouncementManager) CreateAnnounce(announcement *entities.AnnouncementForDB) (int, error) {
	query, args, err := sqlx.In(`INSERT INTO announcements (author_id, author_login, author_email, author_phone, title, description, created_at, photos) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, announcement.AuthorID, announcement.AuthorLogin, announcement.AuthorEmail, announcement.AuthorPhone, announcement.Title, announcement.Description, announcement.CreatedAt, announcement.Photos)
	if err != nil {
		return 0, err
	}

	query = m.db.Rebind(query)
	var id int
	rows, err := m.db.Queryx(query, args...)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}
func (m *AnnouncementManager) UpdateAnnounce() {}
func (m *AnnouncementManager) HideAnnounce()   {}

func (m *AnnouncementManager) DeleteAnnounceById(id string) error {
	query := `DELETE FROM announcements WHERE id=$1`

	_, err := m.db.Query(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("announce with passed id not found")
		}
		return err
	}
	return nil
}

func (m *AnnouncementManager) GetAnnouncePhotosById(id string) ([]string, error) {
	var photos pq.StringArray
	query := `SELECT photos FROM announcements WHERE id=$1`

	row := m.db.QueryRowx(query, id)
	err := row.Scan(&photos)

	return photos, err
}
