package v1

import (
	"log_processor/internal/service"
)

type Handler struct {
	ser *service.LogFilterService
}

func NewHandler(ser *service.LogFilterService) *Handler {
	return &Handler{
		ser: ser,
	}
}
