package handlers

import (
	"github.com/s-aska/Justaway-Ex/crawler/models"
	"github.com/s-aska/anaconda"
)

func HandlerTweet(userIdStr string, data anaconda.Tweet) {
	if data.RetweetedStatus != nil && data.RetweetedStatus.User.IdStr == userIdStr {
		handlerTweetRetweeted(data)
	} else if data.InReplyToUserIdStr == userIdStr {
		handlerTweetReply(data)
	}
}

func handlerTweetRetweeted(data anaconda.Tweet) {
	models.CreateTweetActivityWithReferenceId(
		data.RetweetedStatus.User.IdStr,
		data.RetweetedStatus.IdStr,
		"retweet",
		data.User.IdStr,
		data.IdStr,
		encodeJson(data))
}

func handlerTweetReply(data anaconda.Tweet) {
	models.CreateTweetActivityWithReferenceId(
		data.InReplyToUserIdStr,
		data.InReplyToStatusIdStr,
		"reply",
		data.IdStr,
		data.IdStr,
		encodeJson(data))
}

func HandlerStatusDeletionNotice(data anaconda.StatusDeletionNotice) {
	models.DeleteTweetActivityByStatusId(data.IdStr)
	models.DeleteTweetActivityByReferenceId(data.IdStr)
}
