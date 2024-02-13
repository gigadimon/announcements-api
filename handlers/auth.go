package handlers

import (
	"announce-api/entities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignIn(ctx *gin.Context) {
	user, err := unmarshalInputUser(ctx)
	if err != nil {
		log.Fatalf("unmarshalling input user failed: %s", err.Error())
	}

	h.service.AuthorizeUser(user)
}

func (h *Handler) SignUp(ctx *gin.Context) {
	user, err := unmarshalInputUser(ctx)
	if err != nil {
		defer log.Panicf("unmarshalling input user failed: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unmarshalling input user failed: %s", err.Error())})
		return
	}

	created, err := h.service.CreateUser(user)
	if err != nil {
		defer log.Panic(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(created)
	ctx.JSON(http.StatusOK, gin.H{"message": "user created"})
}

func unmarshalInputUser(ctx *gin.Context) (*entities.InputUser, error) {
	defer ctx.Request.Body.Close()
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}

	user := new(entities.InputUser)

	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return user, nil
}
