package handlers

import (
	"fmt"
	"github.com/s-aska/anaconda"
)

func (h *Handler) HandlerDirectMessage(userId string, data anaconda.DirectMessage) {
	fmt.Printf("[%s] message: @%s => @%s `%s`\n", userId, data.SenderScreenName, data.RecipientScreenName, data.Text)
}

func (h *Handler) HandlerDirectMessageDeletionNotice(userId string, data anaconda.DirectMessageDeletionNotice) {
	fmt.Printf("[%s] message delete: %s:%s\n", userId, data.UserId, data.IdStr)
}
