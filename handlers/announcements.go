package handlers

import (
	"announce-api/entities"
	"announce-api/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (h *Handler) GetAnnouncementList(ctx *gin.Context) {
	id, email, login := ctx.MustGet("id").(float64), ctx.MustGet("email").(string), ctx.MustGet("login").(string)
	fmt.Println(id)
	fmt.Println(email)
	fmt.Println(login)

}

func (h *Handler) UpdateAnnouncement(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(10 << 20)
	if ctx.Request.MultipartForm == nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "missing form data")
		return
	}

	files, ok := ctx.Request.MultipartForm.File["files"]
	if !ok {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "dont have files")
		return
	}

	minioClient, err := minio.New(os.Getenv("MIN_IO_HOST"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MIN_IO_ACCESS_KEY_ID"), os.Getenv("MIN_IO_ACCESS_KEY_SECRET"), ""),
		Secure: false,
	})
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	for _, file := range files {
		fileToUpload, err := file.Open()
		if err != nil {
			utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer fileToUpload.Close()

		ui, err := minioClient.PutObject(context.Background(), "announcements", file.Filename, fileToUpload, file.Size, minio.PutObjectOptions{})
		if err != nil {
			utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("Uploaded", ui.Key, "to", ui.Bucket, ui.ETag, ui.VersionID, ui.Size)
	}
}

func (h *Handler) UploadNewPhoto(ctx *gin.Context) {
}

func (h *Handler) DeletePhoto(ctx *gin.Context) {
}

func (h *Handler) HideAnnounce(ctx *gin.Context) {
}

func (h *Handler) GetAnnouncementById(ctx *gin.Context) {
	postId, _ := ctx.Params.Get("postId")

	announcement, err := h.service.GetOneById(postId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(ctx, http.StatusBadRequest, "announce with passed id not found")
			return
		}

		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": announcement})
}

func (h *Handler) CreateAnnouncement(ctx *gin.Context) {
	inputAnnouncement := new(entities.InputAnnouncement)
	// Парсим форму... !!!Проблема при отправке файлов на 250кбайт+
	form, _ := ctx.MultipartForm()
	if form == nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "missing form data")
		return
	}

	// Если файлы есть - создаем обьекты
	files, ok := form.File["files"]
	photos := make([]string, 0)
	if ok {
		photos = h.service.CreateListOfPhotos(files)
	}

	// Биндим поля формы в структуру, кроме фото, тк требуют отдельной обработки
	if err := ctx.Bind(&inputAnnouncement); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Добавляем фото
	inputAnnouncement.Photos = pq.Array(photos)
	authorInfo := getAuthorInfo(ctx)

	// Наконец создаем анонс
	createdAnnouncement, err := h.service.CreateAnnounce(inputAnnouncement, authorInfo)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "announcement created", "data": createdAnnouncement})
}

func (h *Handler) DeleteAnnouncement(ctx *gin.Context) {
	postId, _ := ctx.Params.Get("postId")

	photos, err := h.service.GetAnnouncePhotosById(postId)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	for _, photo := range photos {
		go h.service.ObjectStorage.DeleteObject(photo)
	}

	if err := h.service.DeleteAnnounceById(postId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(ctx, http.StatusBadRequest, "announce with passed id not found")
			return
		}

		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("announce with id: %s deleted", postId)})
}

func getAuthorInfo(ctx *gin.Context) entities.AuthorInfo {
	authorId, authorEmail, authorLogin := ctx.MustGet("id").(float64), ctx.MustGet("email").(string), ctx.MustGet("login").(string)
	return entities.AuthorInfo{ID: int(authorId), Login: authorLogin, Email: authorEmail}
}
