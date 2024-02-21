package handlers

import (
	"announce-api/entities"
	"announce-api/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) GetGlobalFeed(ctx *gin.Context) {
	page := parseNumericQueryParam(ctx, "page", 1)
	limit := parseNumericQueryParam(ctx, "limit", 20)

	announcementsList, err := h.service.GetGlobalFeed(page, limit)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": announcementsList})
}

func (h *Handler) GetAuthorsAnnouncementList(ctx *gin.Context) {
	page := parseNumericQueryParam(ctx, "page", 1)
	limit := parseNumericQueryParam(ctx, "limit", 20)
	authorId := ctx.MustGet("id").(float64)

	announcementsList, err := h.service.GetAuthorsList(page, limit, int(authorId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("announcements for author id: %s not found", fmt.Sprint(authorId)))
			return
		}
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": announcementsList})
}

func (h *Handler) UpdateAnnouncementById(ctx *gin.Context) {
	inputAnnouncement := new(entities.InputAnnouncement)
	postId, _ := ctx.Params.Get("postId")

	if err := utils.ReadAndUnmarshallInputBody(ctx.Request.Body, inputAnnouncement); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	announcement, err := h.service.UpdateAnnounce(inputAnnouncement, postId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(ctx, http.StatusNotFound, "announcement to update not found")
			return
		}
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": announcement})
}

func (h *Handler) UploadNewAnnouncePhotosById(ctx *gin.Context) {
	postId, _ := ctx.Params.Get("postId")

	form, err := ctx.MultipartForm()
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "parsing form data failed: "+err.Error())
		return
	}

	files, ok := form.File["files"]
	if !ok {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "files required")
		return
	}

	photos := h.service.CreateListOfPhotos(files)

	updatedPhotosList, err := h.service.UploadNewAnnouncePhotosById(pq.Array(photos), postId)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "upload new photos failed. "+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedPhotosList})
}

func (h *Handler) DeleteAnnouncePhotoById(ctx *gin.Context) {
	postId, _ := ctx.Params.Get("postId")
	photoName, _ := ctx.Params.Get("photoName")

	updatedPhotos, err := h.service.DeleteAnnouncePhotoById(postId, photoName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("photo %s doesn't exist in post id %s", photoName, postId))
			return
		}
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedPhotos})
}

func (h *Handler) SwitchAnnounceVisibilityById(ctx *gin.Context) {
	postId, _ := ctx.Params.Get("postId")

	isHidden, err := h.service.SwitchAnnounceVisibilityById(postId)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"is_hidden": isHidden})
}

func (h *Handler) GetAnnouncementById(ctx *gin.Context) {
	postId, _ := ctx.Params.Get("postId")
	userId := ctx.MustGet("id").(float64)

	announcement, err := h.service.GetOneById(postId, fmt.Sprint(userId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(ctx, http.StatusNotFound, "announce with id "+postId+" not found")
			return
		}
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": announcement})
}

func (h *Handler) CreateAnnouncement(ctx *gin.Context) {
	inputAnnouncement := new(entities.InputAnnouncement)
	form, err := ctx.MultipartForm()

	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, "parsing form data failed: "+err.Error())
		return
	}

	files, ok := form.File["files"]
	photos := make([]string, 0)
	if ok {
		photos = h.service.CreateListOfPhotos(files)
	}

	if err := ctx.Bind(&inputAnnouncement); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	authorInfo := getAuthorInfo(ctx)

	createdAnnouncement, err := h.service.CreateAnnounce(inputAnnouncement, pq.Array(photos), authorInfo)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "announcement created", "data": createdAnnouncement})
}

func (h *Handler) DeleteAnnouncementById(ctx *gin.Context) {
	postId, _ := ctx.Params.Get("postId")

	photos, err := h.service.GetAnnouncePhotosById(postId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(ctx, http.StatusNotFound, "announce with id "+postId+" not found")
		}
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	for _, photo := range photos {
		go h.service.ObjectStorage.DeleteObject(os.Getenv("PROJECT_NAME"), photo)
	}

	if err := h.service.DeleteAnnounceById(postId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(ctx, http.StatusNotFound, "announce with id "+postId+" not found")
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

func parseNumericQueryParam(ctx *gin.Context, key string, defaultValue int) (value int) {
	value = defaultValue
	keyParam, ok := ctx.GetQuery(key)
	if ok && keyParam != "" {
		keyInt, err := strconv.Atoi(keyParam)
		if err != nil {
			fmt.Println("error while converting " + key + " string to int")
		} else {
			value = keyInt
		}
	}

	return
}
