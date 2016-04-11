package handlers

import (
	"github.com/s-aska/anaconda"
	"time"
)

func (h *Handler) HandlerEventTweet(userId string, data anaconda.EventTweet) {
	if data.Event.Event == "unfavorite" {
		h.handlerEventTweetUnfavorite(data)
	} else if data.Event.Target.IdStr == userId {
		h.handlerEventTweet(data)
	}
}

func (h *Handler) handlerEventTweet(data anaconda.EventTweet) {
	createdAtTime, err := time.Parse(time.RubyDate, data.Event.CreatedAt)
	if err != nil {
		createdAtTime = time.Now()
	}
	h.model.CreateActivity(
		data.Event.Event,
		data.Event.Target.IdStr,
		data.Event.Source.IdStr,
		data.TargetObject.IdStr,
		createdAtTime.Unix(),
	)
}

func (h *Handler) handlerEventTweetUnfavorite(data anaconda.EventTweet) {
	h.model.DeleteFavoriteActivity(
		data.Event.Source.IdStr,
		data.TargetObject.IdStr)
}
