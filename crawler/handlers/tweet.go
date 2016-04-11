package handlers

import (
	"github.com/s-aska/anaconda"
	"time"
)

func (h *Handler) HandlerTweet(userIdStr string, data anaconda.Tweet) {
	if data.RetweetedStatus != nil && data.RetweetedStatus.User.IdStr == userIdStr {
		h.handlerTweetRetweeted(data)
	} else if data.InReplyToUserIdStr == userIdStr {
		h.handlerTweetReply(data)
	}
}

func (h *Handler) handlerTweetRetweeted(data anaconda.Tweet) {
	createdAtTime, err := data.CreatedAtTime()
	if err != nil {
		createdAtTime = time.Now()
	}
	h.model.CreateRetweetActivity(
		"retweet",
		data.RetweetedStatus.User.IdStr,
		data.User.IdStr,
		data.IdStr,
		data.RetweetedStatus.IdStr,
		createdAtTime.Unix(),
	)
}

func (h *Handler) handlerTweetReply(data anaconda.Tweet) {
	createdAtTime, err := data.CreatedAtTime()
	if err != nil {
		createdAtTime = time.Now()
	}
	h.model.CreateActivity(
		"reply",
		data.InReplyToUserIdStr,
		data.User.IdStr,
		data.IdStr,
		createdAtTime.Unix(),
	)
}

func (h *Handler) HandlerStatusDeletionNotice(data anaconda.StatusDeletionNotice) {
	h.model.DeleteTweetActivity(data.IdStr)
}
