package services

import (
	"announce-api/db"
	"announce-api/entities"
)

type Service struct {
	Auth
	AnnouncementActions
}

type Auth interface {
	CreateUser(user *entities.InputUser) (bool, error)
	AuthorizeUser(user *entities.InputUser)
}

type AnnouncementActions interface {
	GetList()
	GetOneById()
	CreateAnnounce()
	UpdateAnnounce()
	HideAnnounce()
	DeleteAnnounce()
}

func Init(client *db.DatabaseClient) *Service {
	return &Service{
		Auth:                &Authenticator{client},
		AnnouncementActions: &AnnouncementManager{client},
	}
}
