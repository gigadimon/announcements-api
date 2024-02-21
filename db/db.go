package db

import (
	"announce-api/entities"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Config struct {
	Driver   string
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

type DatabaseClient struct {
	Auth
	AnnouncementActions
}

type Auth interface {
	CreateUser(user *entities.InputSignUpUser) (int, error)
	GetUser(user *entities.InputSignInUser) (*entities.User, error)
	UpdateUserToken(user *entities.User, token string) (string, error)
	IsTokenExists(token string) (int, error)
}

type AnnouncementActions interface {
	GetGlobalFeed(page int, limit int) ([]*entities.AnnouncementFromDB, error)
	GetAuthorsList(page int, limit int, authorId int) ([]*entities.AnnouncementFromDB, error)
	GetOneById(postId string, userId string) (*entities.AnnouncementFromDB, error)
	CreateAnnounce(announcement *entities.AnnouncementToDB) (*entities.AnnouncementFromDB, error)
	UpdateAnnounceById(inputAnnouncement *entities.InputAnnouncement, postId string) (*entities.AnnouncementFromDB, error)
	UploadNewAnnouncePhotosById(photos entities.PhotosForDB, id string) (pq.StringArray, error)
	DeleteAnnouncePhotoById(postId string, photoName string) (pq.StringArray, error)
	SwitchAnnounceVisibilityById(postId string) (bool, error)
	DeleteAnnounceById(id string) error
	GetAnnouncePhotosById(id string) ([]string, error)
	GetAnnounceAuthorId(postId string) (string, error)
}

func newClient(db *sqlx.DB) *DatabaseClient {
	return &DatabaseClient{
		Auth:                &Authenticator{db},
		AnnouncementActions: &AnnouncementManager{db},
	}
}

func Connect(config Config) (*DatabaseClient, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", config.Host, config.Port, config.DBName, config.User, config.Password, config.SSLMode)

	db, err := sqlx.Connect(config.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	client := newClient(db)

	db.MustExec(schema)

	fmt.Println("Database connection successfull")
	return client, nil
}
