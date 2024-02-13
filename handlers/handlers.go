package handlers

import "announce-api/services"

type Handler struct {
	service *services.Service
}

func Init(service *services.Service) *Handler {
	return &Handler{service}
}
