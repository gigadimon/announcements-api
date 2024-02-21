package services

import (
	"announce-api/db"
	"announce-api/entities"
	"errors"

	"github.com/lib/pq"
)

type AnnouncementManager struct {
	client *db.DatabaseClient
}

func (m *AnnouncementManager) GetGlobalFeed(page int, limit int) ([]*entities.AnnouncementFromDB, error) {
	return m.client.GetGlobalFeed(page, limit)
}

func (m *AnnouncementManager) GetAuthorsList(page int, limit int, authorId int) ([]*entities.AnnouncementFromDB, error) {
	return m.client.GetAuthorsList(page, limit, authorId)
}

func (m *AnnouncementManager) GetOneById(postId string, userId string) (*entities.AnnouncementFromDB, error) {
	return m.client.GetOneById(postId, userId)
}

func (m *AnnouncementManager) CreateAnnounce(inputAnnouncement *entities.InputAnnouncement, photos entities.PhotosForDB, author entities.AuthorInfo) (*entities.AnnouncementFromDB, error) {
	if err := inputAnnouncement.Validate(); err != nil {
		return new(entities.AnnouncementFromDB), errors.New("passed announcement incorrect: " + err.Error())
	}

	announcement := entities.NewAnnouncementForDB(inputAnnouncement, photos, author)

	announcementFromDB, err := m.client.CreateAnnounce(announcement)
	if err != nil {
		return new(entities.AnnouncementFromDB), err
	}

	return announcementFromDB, nil
}

func (m *AnnouncementManager) UpdateAnnounce(inputAnnouncement *entities.InputAnnouncement, postId string) (*entities.AnnouncementFromDB, error) {
	if err := inputAnnouncement.Validate(); err != nil {
		return new(entities.AnnouncementFromDB), errors.New("passed announcement incorrect: " + err.Error())
	}

	return m.client.UpdateAnnounceById(inputAnnouncement, postId)
}

func (m *AnnouncementManager) SwitchAnnounceVisibilityById(postId string) (bool, error) {
	return m.client.SwitchAnnounceVisibilityById(postId)
}

func (m *AnnouncementManager) GetAnnouncePhotosById(id string) ([]string, error) {
	return m.client.GetAnnouncePhotosById(id)
}

func (m *AnnouncementManager) DeleteAnnounceById(id string) error {
	return m.client.DeleteAnnounceById(id)
}

func (m *AnnouncementManager) UploadNewAnnouncePhotosById(photos entities.PhotosForDB, id string) (pq.StringArray, error) {
	return m.client.UploadNewAnnouncePhotosById(photos, id)
}

func (m *AnnouncementManager) DeleteAnnouncePhotoById(postId string, photoName string) (pq.StringArray, error) {
	return m.client.DeleteAnnouncePhotoById(postId, photoName)
}

func (m *AnnouncementManager) IsUserAnnounceAuthor(postId string, userId string) (bool, error) {
	authorId, err := m.client.GetAnnounceAuthorId(postId)
	if err != nil {
		return false, err
	}

	if authorId != userId {
		return false, errors.New("you are not announcement author")
	}

	return true, nil
}
