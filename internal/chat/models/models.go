package models

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ChatStatus uint8

var BotAPI *tgbotapi.BotAPI

const (
	MainStatus = ChatStatus(iota)
	NilStatus  = ChatStatus(iota)
)
