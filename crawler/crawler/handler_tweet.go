package crawler

import (
	"github.com/s-aska/anaconda"
)

func handlerTweet(userIdStr string, data anaconda.Tweet) {
	if data.RetweetedStatus != nil && data.RetweetedStatus.User.IdStr == userIdStr {
		handlerTweetRetweeted(data)
	} else if data.InReplyToUserIdStr == userIdStr {
		handlerTweetReply(data)
	}
}

func handlerTweetRetweeted(data anaconda.Tweet) {
	createActivityWithReferenceId(
		data.RetweetedStatus.User.IdStr,
		data.RetweetedStatus.IdStr,
		"retweet",
		data.User.IdStr,
		data.IdStr,
		encodeJson(data))
}

func handlerTweetReply(data anaconda.Tweet) {
	createActivityWithReferenceId(
		data.InReplyToUserIdStr,
		data.InReplyToStatusIdStr,
		"reply",
		data.IdStr,
		data.IdStr,
		encodeJson(data))
}

func handlerStatusDeletionNotice(data anaconda.StatusDeletionNotice) {
	deleteActivityByStatusId(data.IdStr)
	deleteActivityByReferenceId(data.IdStr)
}
