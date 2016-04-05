package handlers

import (
	"github.com/s-aska/Justaway-Ex/crawler/models"
	"github.com/s-aska/anaconda"
)

func HandlerEventTweet(userId string, data anaconda.EventTweet) {
	if data.Event.Event == "quoted_tweet" && data.TargetObject.QuotedStatus.User.IdStr == userId {
		handlerEventTweetQuotedTweet(data)
	} else if data.TargetObject.User.IdStr == userId {
		switch data.Event.Event {
		case "favorite", "favorited_retweet", "retweeted_retweet":
			handlerEventTweetFavorite(data)
		case "unfavorite":
			handlerEventTweetUnfavorite(data)
		}
	}
}

func handlerEventTweetQuotedTweet(data anaconda.EventTweet) {
	models.CreateTweetActivityWithReferenceId(
		data.TargetObject.QuotedStatus.User.IdStr,
		data.TargetObject.QuotedStatus.IdStr,
		data.Event.Event,
		data.Event.Source.IdStr,
		data.TargetObject.IdStr,
		encodeJson(data))
}

func handlerEventTweetFavorite(data anaconda.EventTweet) {
	models.CreateTweetActivity(
		data.TargetObject.User.IdStr,
		data.TargetObject.IdStr,
		data.Event.Event,
		data.Event.Source.IdStr,
		encodeJson(data))
}

func handlerEventTweetUnfavorite(data anaconda.EventTweet) {
	models.DeleteTweetActivity(
		data.TargetObject.IdStr,
		"favorite",
		data.Event.Source.IdStr)
}
