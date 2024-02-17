package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func SendErrorResponse(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"error": message})
	ctx.AbortWithError(status, errors.New(message))
}
