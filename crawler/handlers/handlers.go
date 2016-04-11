package handlers

import (
	"github.com/s-aska/Justaway-Ex/crawler/models"
)

type Handler struct {
	model *models.Model
}

func New(model *models.Model) *Handler {
	return &Handler{
		model: model,
	}
}
