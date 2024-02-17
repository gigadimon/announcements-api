package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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
		Creds:  credentials.NewStaticV4(os.Getenv("MIN_IO_ACCESS_KEY_ID"), os.Getenv("MIN_IO_ACCESS_KEY_SECRET"), ""),
		Secure: false,
	})
}

func (m *Minio) CreateObject(bucketName string, fileExt string, fileToUpload multipart.File, size int64) (minio.UploadInfo, error) {
	return m.client.PutObject(context.Background(), bucketName, generateFilename(fileExt), fileToUpload, size, minio.PutObjectOptions{})
}

func (m *Minio) DeleteObject(photoPath string) error {
	photoParts, err := parsePhotoPath(photoPath)
	if err != nil {
		return err
	}
	bucketName, photoName := photoParts[0], photoParts[1]
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

	photos <- fmt.Sprintf("/%s/%s", bucketName, ui.Key)
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

func parsePhotoPath(path string) ([2]string, error) {
	photoParts := strings.Split(path, "/")
	if len(photoParts) != 3 {
		return [2]string{}, errors.New("photo path is incorrect")
	}

	return [2]string{photoParts[1], photoParts[2]}, nil
}
