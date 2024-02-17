package services

import (
	"announce-api/db"
	"announce-api/entities"
	"mime/multipart"
	"sync"

	"github.com/minio/minio-go/v7"
)

type Service struct {
	Auth
	AnnouncementActions
	ObjectStorage
}

type Auth interface {
	CreateUser(user *entities.InputSignUpUser) (int, error)
	AuthorizeUser(user *entities.InputSignInUser) (string, error)
	IsTokenExists(token string) (int, error)
}

type AnnouncementActions interface {
	GetList()
	GetOneById(id string) (*entities.AnnouncementFromDB, error)
	CreateAnnounce(announcement *entities.InputAnnouncement, author entities.AuthorInfo) (*entities.AnnouncementForDB, error)
	UpdateAnnounce()
	HideAnnounce()
	DeleteAnnounceById(id string) error
	GetAnnouncePhotosById(id string) ([]string, error)
}

type ObjectStorage interface {
	CreateObject(bucketName string, fileExt string, fileToUpload multipart.File, size int64) (minio.UploadInfo, error)
	DeleteObject(photoPath string) error
	СreatePhotoObjectFromFile(bucketName string, file *multipart.FileHeader, photos chan string, wg *sync.WaitGroup)
	CreateListOfPhotos(files []*multipart.FileHeader) (photos []string)
}

type InitParams struct {
	DBClient            *db.DatabaseClient
	ObjectStorageClient *minio.Client
}

func Init(params InitParams) *Service {
	return &Service{
		Auth:                &Authenticator{client: params.DBClient},
		AnnouncementActions: &AnnouncementManager{client: params.DBClient},
		ObjectStorage:       &Minio{client: params.ObjectStorageClient},
	}
}
