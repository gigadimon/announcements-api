package services

import (
	"announce-api/db"
	"announce-api/entities"
	"errors"
	"time"
)

type AnnouncementManager struct {
	client *db.DatabaseClient
}

func (m *AnnouncementManager) GetList() {}
func (m *AnnouncementManager) GetOneById(id string) (*entities.AnnouncementFromDB, error) {
	return m.client.GetOneById(id)
}

func (m *AnnouncementManager) CreateAnnounce(inputAnnouncement *entities.InputAnnouncement, author entities.AuthorInfo) (*entities.AnnouncementForDB, error) {
	if err := inputAnnouncement.Validate(); err != nil {
		return new(entities.AnnouncementForDB), errors.New("passed announcement incorrect: " + err.Error())
	}

	announcement := &entities.AnnouncementForDB{
		AuthorID:    author.ID,
		AuthorEmail: author.Email,
		AuthorLogin: author.Login,
		AuthorPhone: inputAnnouncement.AuthorPhone,
		Photos:      inputAnnouncement.Photos,
		Title:       inputAnnouncement.Title,
		Description: inputAnnouncement.Description,
		CreatedAt:   time.Now(),
	}

	announcementId, err := m.client.CreateAnnounce(announcement)
	if err != nil {
		return new(entities.AnnouncementForDB), err
	}

	announcement.ID = announcementId
	return announcement, nil
}

func (m *AnnouncementManager) UpdateAnnounce() {}
func (m *AnnouncementManager) HideAnnounce()   {}

func (m *AnnouncementManager) GetAnnouncePhotosById(id string) ([]string, error) {
	return m.client.GetAnnouncePhotosById(id)
}
func (m *AnnouncementManager) DeleteAnnounceById(id string) error {
	return m.client.DeleteAnnounceById(id)
}
