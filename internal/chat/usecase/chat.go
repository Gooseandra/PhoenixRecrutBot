package usecase

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"phenixRecrutBot/internal/chat/models"
	"phenixRecrutBot/internal/constants"
	"phenixRecrutBot/internal/pkg/IAP/IAP"
	"phenixRecrutBot/internal/pkg/register"
	"strconv"
	"sync"
	"time"
)

type Chat struct {
	Id         int64
	Channel    chan tgbotapi.Update
	TimeStart  time.Time
	TimeFinish time.Time
	Iap        IAP.IAP
	Register   register.RegisterUC
}

func (chat Chat) Routine(chats map[int64]Chat, mainMutex *sync.Mutex) {
	lastMassageTime := time.After(time.Minute * 10)
	status := models.NilStatus
	for {
		select {
		case message := <-chat.Channel:
			lastMassageTime = time.After(time.Hour * 10)
			log.Println("MainStatus.Update")
			if message.Message != nil {
				switch status {
				case models.NilStatus:
					switch message.Message.Text {
					case "/start":
						status = models.MainStatus
						showCmd := tgbotapi.NewMessage(chat.Id, constants.HelloMessage)
						showCmd.ReplyMarkup = constants.StartKeyboard
						models.BotAPI.Request(showCmd)
					}
				case models.MainStatus:
					switch message.Message.Text {
					case constants.AdminCodeMessage:
						msg := tgbotapi.NewMessage(chat.Id, "Авторизация")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						models.BotAPI.Request(msg)
						chat.Register.GetForms(chat.Id, chat.Channel)
					default:
						msg := tgbotapi.NewMessage(chat.Id, "")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						models.BotAPI.Request(msg)
						_, err := chat.Register.Register(chat.Id, chat.Channel, message.Message.From.UserName)
						if err != nil {
							log.Println(err.Error())
						}
					}
				}
			}

		case <-lastMassageTime:
			log.Println("time out")
			mainMutex.Lock()
			log.Println("Chat " + strconv.FormatInt(chat.Id, 10) + " deleted")
			log.Println(chats)
			delete(chats, chat.Id)
			mainMutex.Unlock()
		}
	}
}
