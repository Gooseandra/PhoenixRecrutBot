package register

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"phenixRecrutBot/internal/pkg/register/models"
)

type RegisterUC interface {
	Register(id int64, channel chan tgbotapi.Update, userName string) (bool, error)
	GetForms(id int64, channel chan tgbotapi.Update) error
}

type RegisterRepository interface {
	SetNewForm(form models.Form, id int64) error
	List() ([]models.Form, error)
}
