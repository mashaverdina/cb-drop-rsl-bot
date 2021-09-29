package chatutils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message interface {
	Chat() int64
}

type MessageFromChat interface {
	Message
	MessageID() int
}

type MessageWithOldVersion interface {
	MessageFromChat
	Original() string
}

func EditTo(msg MessageFromChat, text string, markup *tgbotapi.InlineKeyboardMarkup) []tgbotapi.Chattable {
	resp := tgbotapi.NewEditMessageText(msg.Chat(), msg.MessageID(), text)
	if markup != nil {
		resp.ReplyMarkup = markup
	}
	resp.ParseMode = tgbotapi.ModeMarkdown
	return []tgbotapi.Chattable{resp}
}

func TextTo(msg Message, text string, markup interface{}) []tgbotapi.Chattable {
	resp := tgbotapi.NewMessage(msg.Chat(), text)
	if markup != nil {
		resp.ReplyMarkup = markup
	}
	resp.ParseMode = tgbotapi.ModeMarkdown
	return []tgbotapi.Chattable{resp}
}

func DisableKeyboardAndSendNew(msg MessageWithOldVersion, text string, markup interface{}) []tgbotapi.Chattable {
	return JoinResp(
		EditTo(msg, msg.Original(), nil),
		TextTo(msg, text, markup),
	)
}

func RemoveAndSendNew(msg MessageFromChat, text string, markup interface{}) []tgbotapi.Chattable {
	return JoinResp(
		[]tgbotapi.Chattable{tgbotapi.NewDeleteMessage(msg.Chat(), msg.MessageID())},
		TextTo(msg, text, markup),
	)
}

func JoinResp(resps ...[]tgbotapi.Chattable) []tgbotapi.Chattable {
	result := make([]tgbotapi.Chattable, 0)
	for _, arr := range resps {
		result = append(result, arr...)
	}
	return result
}
