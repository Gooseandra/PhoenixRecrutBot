package IAP

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"phenixRecrutBot/internal/chat/models"
	"phenixRecrutBot/internal/constants"
	md "phenixRecrutBot/internal/pkg/IAP"
)

type IAP struct {
}

func New() IAP {
	return IAP{}
}

func (iap IAP) WrongInput(id int64) {
	models.BotAPI.Send(tgbotapi.NewMessage(id, constants.WrongInputErr))
}

func (iap IAP) PrintText(id int64, text string) {
	models.BotAPI.Send(tgbotapi.NewMessage(id, text))
}

func (iap IAP) InputText(id int64, channel chan tgbotapi.Update, disc string) (string, *md.MetaData, error) {
	models.BotAPI.Send(tgbotapi.NewMessage(id, disc))
	message := <-channel
	if message.Message == nil {
		return "", nil, errors.New(constants.NilInputErr)
	}
	metadata := md.MetaData{
		UserName: message.Message.From.UserName,
		Name:     message.Message.From.FirstName,
		Surname:  message.Message.From.LastName}
	return message.Message.Text, &metadata, nil
}
