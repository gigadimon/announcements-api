package services

import (
	"announce-api/db"
	"announce-api/entities"
	"mime/multipart"
	"sync"

	"github.com/lib/pq"
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
	GetGlobalFeed(page int, limit int) ([]*entities.AnnouncementFromDB, error)
	GetAuthorsList(page int, limit int, authorId int) ([]*entities.AnnouncementFromDB, error)
	GetOneById(postId string, userId string) (*entities.AnnouncementFromDB, error)
	CreateAnnounce(inputAnnouncement *entities.InputAnnouncement, photos entities.PhotosForDB, author entities.AuthorInfo) (*entities.AnnouncementFromDB, error)
	UpdateAnnounce(inputAnnouncement *entities.InputAnnouncement, postId string) (*entities.AnnouncementFromDB, error)
	UploadNewAnnouncePhotosById(photos entities.PhotosForDB, postId string) (pq.StringArray, error)
	DeleteAnnouncePhotoById(postId string, photoName string) (pq.StringArray, error)
	SwitchAnnounceVisibilityById(postId string) (bool, error)
	DeleteAnnounceById(postId string) error
	GetAnnouncePhotosById(postId string) ([]string, error)
	IsUserAnnounceAuthor(postId string, userId string) (bool, error)
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
