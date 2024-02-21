package services

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	client *minio.Client
}

const (
	defaultBucketName = "announcements"
)

func NewMinioStorage() (*minio.Client, error) {
	return minio.New(os.Getenv("MIN_IO_HOST"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ROOT_USER"), os.Getenv("MINIO_ROOT_PASSWORD"), ""),
		Secure: false,
	})
}

func (m *Minio) CreateObject(bucketName string, fileExt string, fileToUpload multipart.File, size int64) (minio.UploadInfo, error) {
	return m.client.PutObject(context.Background(), bucketName, generateFilename(fileExt), fileToUpload, size, minio.PutObjectOptions{})
}

func (m *Minio) DeleteObject(bucketName string, photoName string) error {
	if err := m.client.RemoveObject(context.Background(), bucketName, photoName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}

func (m *Minio) СreatePhotoObjectFromFile(bucketName string, file *multipart.FileHeader, photos chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileToUpload, err := file.Open()
	if err != nil {
		photos <- "" // ?? какая практика вообще в этом случае и что делать с рутинами, которые зафейлились? решил ничего не делать потому, что не хочу отвечать фейлом, если зафейлился какой-то один файл, но, мб, как-то уведомлять юзера все-таки в респонсе...
		return
	}
	defer fileToUpload.Close()

	ui, err := m.CreateObject(bucketName, filepath.Ext(file.Filename), fileToUpload, file.Size)
	if err != nil {
		photos <- "" // ??
		return
	}

	photos <- ui.Key
}

func (m *Minio) CreateListOfPhotos(files []*multipart.FileHeader) (photos []string) {
	photosChannel := make(chan string, len(files))
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go m.СreatePhotoObjectFromFile(defaultBucketName, file, photosChannel, &wg)
	}
	wg.Wait()
	close(photosChannel)

	for v := range photosChannel {
		if v != "" {
			photos = append(photos, v)
		}
	}

	return photos
}

func generateFilename(fileExt string) string {
	return uuid.NewString() + fileExt
}
