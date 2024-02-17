package handlers

import (
	"announce-api/utils"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) Authenticate(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")
	tokenPart, err := getTokenFromAuthHeader(authHeader)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = h.service.IsTokenExists(tokenPart)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "not authorized: "+err.Error())
		return
	}

	token, err := parseToken(tokenPart)
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "not authorized: token invalid. "+err.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "not authorized: token invalid")
		return
	}

	expires, ok := claims["expires"].(float64)
	if !ok {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "not authorized: token invalid. expires field is incorrect")
		return
	}

	expiresDiffWithNow := int64(expires) - time.Now().Unix()
	if expiresDiffWithNow <= 0 {
		utils.SendErrorResponse(ctx, http.StatusUnauthorized, "not authorized: token expired")
		return
	}

	ctx.Set("id", claims["id"])
	ctx.Set("email", claims["email"])
	ctx.Set("login", claims["login"])
	ctx.Next()
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	return token, err
}

func getTokenFromAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("not authorized: authorization header required")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return "", errors.New("not authorized: token invalid")
	}

	tokenPart := authHeaderParts[1]
	if len(tokenPart) == 0 {
		return "", errors.New("not authorized: token invalid")
	}

	return tokenPart, nil
}

func (h *Handler) IsUserAnnounceAuthor(ctx *gin.Context) {
	postId, _ := ctx.Params.Get("postId")
	userId := ctx.MustGet("id").(float64)

	_, err := h.service.IsUserAnnounceAuthor(postId, fmt.Sprint(int(userId)))
	if err != nil {
		utils.SendErrorResponse(ctx, http.StatusForbidden, err.Error())
		return
	}

	ctx.Next()
}
