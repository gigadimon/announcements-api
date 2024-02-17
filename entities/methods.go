package entities

import (
	"announce-api/utils"
	"time"
)

// Сделал эти методы для того, чтобы в ValidateStruct использовать структуры, которые соответствуют интерфейсу Validate, но, шото, идея не очень получилась)
func (u *InputSignInUser) Validate() error {
	if err := utils.ValidateStruct(u); err != nil {
		return err
	}
	return nil
}

func (u *InputSignUpUser) Validate() error {
	if err := utils.ValidateStruct(u); err != nil {
		return err
	}
	return nil
}

func (a *InputAnnouncement) Validate() error {
	if err := utils.ValidateStruct(a); err != nil {
		return err
	}
	return nil
}

func (a *AnnouncementToDB) Validate() error {
	if err := utils.ValidateStruct(a); err != nil {
		return err
	}
	return nil
}

func NewAnnouncementForDB(inputAnnouncement *InputAnnouncement, photos PhotosForDB, author AuthorInfo) *AnnouncementToDB {
	return &AnnouncementToDB{
		AuthorID:    author.ID,
		AuthorEmail: author.Email,
		AuthorLogin: author.Login,
		AuthorPhone: inputAnnouncement.AuthorPhone,
		Photos:      photos,
		Title:       inputAnnouncement.Title,
		Description: inputAnnouncement.Description,
		CreatedAt:   time.Now(),
	}
}
