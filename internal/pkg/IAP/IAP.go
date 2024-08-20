package IAP

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MetaData struct {
	Name     string
	Surname  string
	UserName string
}

type IAP interface {
	WrongInput(id int64)
	PrintText(id int64, text string)
	InputText(id int64, channel chan tgbotapi.Update, disc string) (string, *MetaData, error)
}
