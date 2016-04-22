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
	h.enqueuer.Enqueue(
		"resque:queue:default",
		"NotificationTweet",
		data.RetweetedStatus.User.IdStr,
		data.User.ScreenName,
		"retweet",
		data.RetweetedStatus.Text,
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
	h.enqueuer.Enqueue(
		"resque:queue:default",
		"NotificationTweet",
		data.InReplyToUserIdStr,
		data.User.ScreenName,
		"reply",
		data.Text,
	)
}

func (h *Handler) HandlerStatusDeletionNotice(data anaconda.StatusDeletionNotice) {
	h.model.DeleteTweetActivity(data.IdStr)
}
