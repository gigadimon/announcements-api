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

func (h *Handler) Authenticate(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	tokenPart, err := getTokenFromAuthHeader(authHeader)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	_, err = h.service.IsTokenExists(tokenPart)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusForbidden, "not authorized: "+err.Error())
		return
	}

	token, err := parseToken(tokenPart)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusForbidden, "not authorized: token invalid. "+err.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.SendErrorResponse(c, http.StatusForbidden, "not authorized: token invalid")
		return
	}

	expires, ok := claims["expires"].(float64)
	if !ok {
		utils.SendErrorResponse(c, http.StatusForbidden, "not authorized: token invalid. expires field is incorrect")
		return
	}

	expiresDiffWithNow := int64(expires) - time.Now().Unix()
	if expiresDiffWithNow <= 0 {
		utils.SendErrorResponse(c, http.StatusForbidden, "not authorized: token expired")
		return
	}

	c.Set("id", claims["id"])
	c.Set("email", claims["email"])
	c.Set("login", claims["login"])
	c.Next()
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
