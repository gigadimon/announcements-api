package handlers

import (
	"announce-api/entities"
	"announce-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignIn(ctx *gin.Context) {
	defer ctx.Request.Body.Close()
	user := new(entities.InputSignInUser)
	if err := utils.ReadAndUnmarshallInputBody(ctx.Request.Body, user); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.AuthorizeUser(user)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) SignUp(ctx *gin.Context) {
	defer ctx.Request.Body.Close()
	user := new(entities.InputSignUpUser)
	if err := utils.ReadAndUnmarshallInputBody(ctx.Request.Body, user); err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateUser(user)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user created", "userId": id})
}
