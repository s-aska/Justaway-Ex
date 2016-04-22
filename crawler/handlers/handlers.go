package handlers

import (
	"github.com/kavu/go-resque"
	_ "github.com/kavu/go-resque/godis" // Use godis driver
	"github.com/s-aska/Justaway-Ex/crawler/models"
	"github.com/simonz05/godis/redis"
)

type Handler struct {
	model    *models.Model
	enqueuer *resque.RedisEnqueuer
}

func New(model *models.Model) *Handler {
	client := redis.New("tcp:127.0.0.1:6379", 0, "")
	enqueuer := resque.NewRedisEnqueuer("godis", client)
	return &Handler{
		model:    model,
		enqueuer: enqueuer,
	}
}
