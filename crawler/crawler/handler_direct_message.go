package crawler

import (
	"fmt"
	"github.com/s-aska/anaconda"
)

func handlerDirectMessage(userId string, data anaconda.DirectMessage) {
	fmt.Printf("[%s] message: @%s => @%s `%s`\n", userId, data.SenderScreenName, data.RecipientScreenName, data.Text)
}

func handlerDirectMessageDeletionNotice(userId string, data anaconda.DirectMessageDeletionNotice) {
	fmt.Printf("[%s] message delete: %s:%s\n", userId, data.UserId, data.IdStr)
}
