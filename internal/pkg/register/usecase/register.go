package usecase

import (
	errors2 "errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	models2 "phenixRecrutBot/internal/chat/models"
	"phenixRecrutBot/internal/constants"
	"phenixRecrutBot/internal/pkg/IAP/IAP"
	"phenixRecrutBot/internal/pkg/errors"
	"phenixRecrutBot/internal/pkg/register/models"
	"phenixRecrutBot/internal/pkg/register/repository"
	"strconv"
	"strings"
)

type RegisterUseCase struct {
	iap  IAP.IAP
	repo repository.RegisterRepo
}

func NewRegisterUseCase(iap IAP.IAP, repo repository.RegisterRepo) RegisterUseCase {
	return RegisterUseCase{iap: iap, repo: repo}
}

func (r RegisterUseCase) Register(id int64, channel chan tgbotapi.Update, message *tgbotapi.Message) (bool, error) {
	forms, err := r.repo.List()
	if err != nil {
		return false, err
	}
	if message.Text != constants.Retry {
		for _, item := range forms {
			if item.Tg == message.From.UserName {
				r.iap.PrintText(id, "Ты уже создал анкету")
				return false, nil
			}
		}
	}
	name, metadata, err := r.iap.InputText(id, channel, constants.ReqNameMessage)
	if err != nil {
		errors.Warn(id, err.Error())
	}
	confirmed := false
	fullName := strings.Split(name, " ")
	if len(fullName) < 2 {
		errors.Warn(id, "can`t parse name")
	} else {
		if strings.ToLower(fullName[0]) == strings.ToLower(metadata.Name) &&
			strings.ToLower(fullName[1]) == strings.ToLower(metadata.Surname) {
			confirmed = true
		}
	}

	phone, _, err := r.iap.InputText(id, channel, constants.ReqPhoneMessage)
	if err != nil {
		log.Println("WARNING: %s", err.Error())
	}

	vk, _, err := r.iap.InputText(id, channel, constants.ReqVkMessage)
	if err != nil {
		log.Println("WARNING: %s", err.Error())
	}

	if len(fullName) == 1 {
		fullName = append(fullName, "")
	} else if len(fullName) == 0 {
		fullName = append(fullName, "")
		fullName = append(fullName, "")
	}

	form := models.Form{
		ConfirmedName: confirmed,
		Name:          fullName[0],
		Surname:       fullName[1],
		Vk:            vk,
		Tg:            metadata.UserName,
		Phone:         phone}

	err = r.repo.SetNewForm(form, id)
	if err != nil {
		return false, err
	}

	msg := tgbotapi.NewMessage(id, constants.FinishFormMessage)
	msg.ReplyMarkup = constants.RetryKeyboard
	_, err = models2.BotAPI.Request(msg)
	if err != nil {
		return false, err
	}
	return true, err
}

func (r RegisterUseCase) GetForms(id int64, channel chan tgbotapi.Update) error {
	_, metadata, err := r.iap.InputText(id, channel, "")
	if err != nil {
		errors.Warn(id, err.Error())
	}
	if metadata.UserName != "Gooseandra" {
		errors.Warn(id, "non-admin enter from "+metadata.UserName+" "+metadata.Name+" "+metadata.Surname)
		return errors2.New("non-admin enter")
	}

	forms, err := r.repo.List()
	if err != nil {
		errors.Warn(id, err.Error())
		return err
	}

	res := "Количество заявок: " + strconv.Itoa(len(forms)) + "\n\n"

	for _, item := range forms {
		res += "Имя: " + item.Name + " " + item.Surname + "\n" +
			"tg: " + item.Tg + "\n" +
			"Телефон: " + item.Phone
		if item.Vk != "0" {
			res += "\n" + "vk: " + item.Vk
		}
		res += "\n\n"
	}

	r.iap.PrintText(id, res)
	return nil
}
