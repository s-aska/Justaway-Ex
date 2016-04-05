package crawler

import (
	"github.com/s-aska/anaconda"
)

func handlerEventTweet(userId string, data anaconda.EventTweet) {
	if data.Event.Event == "quoted_tweet" && data.TargetObject.QuotedStatus.User.IdStr == userId {
		handlerEventTweetQuotedTweet(data)
	} else if data.TargetObject.User.IdStr == userId {
		if data.Event.Event == "favorite" {
			handlerEventTweetFavorite(data)
		} else if data.Event.Event == "favorited_retweet" {
			handlerEventTweetFavorite(data)
		} else if data.Event.Event == "retweeted_retweet" {
			handlerEventTweetFavorite(data)
		} else if data.Event.Event == "unfavorite" {
			handlerEventTweetUnfavorite(data)
		}
	}
}

func handlerEventTweetQuotedTweet(data anaconda.EventTweet) {
	createActivityWithReferenceId(
		data.TargetObject.QuotedStatus.User.IdStr,
		data.TargetObject.QuotedStatus.IdStr,
		data.Event.Event,
		data.Event.Source.IdStr,
		data.TargetObject.IdStr,
		encodeJson(data))
}

func handlerEventTweetFavorite(data anaconda.EventTweet) {
	createActivity(
		data.TargetObject.User.IdStr,
		data.TargetObject.IdStr,
		data.Event.Event,
		data.Event.Source.IdStr,
		encodeJson(data))
}

func handlerEventTweetUnfavorite(data anaconda.EventTweet) {
	deleteActivity(
		data.TargetObject.IdStr,
		"favorite",
		data.Event.Source.IdStr)
}
