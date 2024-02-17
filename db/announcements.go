package db

import (
	"announce-api/entities"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AnnouncementManager struct {
	db *sqlx.DB
}

func (m *AnnouncementManager) GetList(page int, limit int) ([]*entities.AnnouncementFromDB, error) {
	announcementsList := make([]*entities.AnnouncementFromDB, 0)
	query := `SELECT * FROM announcements WHERE is_hidden=false ORDER BY created_at DESC OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`

	return announcementsList, m.db.Select(&announcementsList, query, (page-1)*limit, limit)
}

func (m *AnnouncementManager) SwitchAnnounceVisibilityById(postId string) (bool, error) {
	var isHidden bool
	query := `UPDATE announcements SET is_hidden=CASE WHEN is_hidden=TRUE THEN FALSE ELSE TRUE END WHERE id=$1 RETURNING is_hidden`

	row := m.db.QueryRow(query, postId)
	if err := row.Scan(&isHidden); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return isHidden, errors.New("announce with passed id not found")
		}
		return isHidden, err
	}

	return isHidden, nil
}

func (m *AnnouncementManager) GetOneById(postId string, userId string) (*entities.AnnouncementFromDB, error) {
	announcement := new(entities.AnnouncementFromDB)
	query := `SELECT * FROM announcements WHERE id=$1 AND (is_hidden=FALSE OR (is_hidden=TRUE AND author_id=$2))`

	row := m.db.QueryRowx(query, postId, userId)
	if err := row.StructScan(announcement); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(entities.AnnouncementFromDB), errors.New("announce with id " + postId + " not found")
		}
		return new(entities.AnnouncementFromDB), err
	}
	return announcement, nil
}

func (m *AnnouncementManager) CreateAnnounce(announcement *entities.AnnouncementToDB) (*entities.AnnouncementFromDB, error) {
	announcementFromDB := new(entities.AnnouncementFromDB)
	query, args, err := sqlx.In(`INSERT INTO announcements (author_id, author_login, author_email, author_phone, title, description, created_at, photos) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`, announcement.AuthorID, announcement.AuthorLogin, announcement.AuthorEmail, announcement.AuthorPhone, announcement.Title, announcement.Description, announcement.CreatedAt, announcement.Photos)
	if err != nil {
		return announcementFromDB, err
	}

	query = m.db.Rebind(query)
	row := m.db.QueryRowx(query, args...)
	if err := row.StructScan(announcementFromDB); err != nil {
		return announcementFromDB, err
	}

	return announcementFromDB, nil
}

func (m *AnnouncementManager) UpdateAnnounceById(inputAnnouncement *entities.InputAnnouncement, postId string) (*entities.AnnouncementFromDB, error) {
	query := `UPDATE announcements SET author_phone=$1, title=$2, description=$3 WHERE id=$4 RETURNING *`

	announcement := new(entities.AnnouncementFromDB)
	row := m.db.QueryRowx(query, inputAnnouncement.AuthorPhone, inputAnnouncement.Title, inputAnnouncement.Description, postId)

	if err := row.StructScan(announcement); err != nil {
		return new(entities.AnnouncementFromDB), err
	}

	return announcement, nil
}

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

func (m *AnnouncementManager) UploadNewAnnouncePhotosById(photos entities.PhotosForDB, postId string) (pq.StringArray, error) {
	var updatedPhotos pq.StringArray
	query := `UPDATE announcements SET photos=array_cat(photos, $1) WHERE id=$2 RETURNING photos`

	row := m.db.QueryRowx(query, photos, postId)
	if err := row.Scan(&updatedPhotos); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return updatedPhotos, errors.New("announce id " + postId + " doesn't exists")
		}

		return updatedPhotos, err
	}

	return updatedPhotos, nil
}

func (m *AnnouncementManager) DeleteAnnouncePhotoById(postId string, photoName string) (pq.StringArray, error) {
	var updatedPhotos pq.StringArray
	query := `UPDATE announcements SET photos=array_remove(photos, $1) WHERE id=$2 RETURNING photos`

	photoPath := fmt.Sprint("/announcements/" + photoName)
	row := m.db.QueryRow(query, photoPath, postId)
	if err := row.Scan(&updatedPhotos); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return updatedPhotos, errors.New("announce id " + postId + " doesn't exists")
		}

		return updatedPhotos, err
	}

	return updatedPhotos, nil
}

func (m *AnnouncementManager) GetAnnounceAuthorId(postId string) (string, error) {
	var id string
	query := `SELECT author_id FROM announcements WHERE id=$1`

	row := m.db.QueryRow(query, postId)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("announce id " + postId + " doesn't exists")
		}
		return "", err
	}

	return id, nil

}
