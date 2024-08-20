package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/go-yaml/yaml"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"phenixRecrutBot/config"
	"phenixRecrutBot/internal/chat/models"
	"phenixRecrutBot/internal/chat/usecase"
	"phenixRecrutBot/internal/pkg/IAP/IAP"
	"phenixRecrutBot/internal/pkg/register/repository"
	usecase2 "phenixRecrutBot/internal/pkg/register/usecase"
	"strconv"
	"sync"
)

func main() {
	var mainMutex sync.Mutex
	var chats = map[int64]usecase.Chat{}
	var settings config.Settings
	bytes, fail := os.ReadFile("./config/config.yml")
	if fail != nil {
		log.Panic(fail.Error())
	}
	fail = yaml.Unmarshal([]byte(bytes), &settings)
	if fail != nil {
		log.Panic(fail.Error())
	}
	log.Println(settings)
	models.BotAPI, fail = tgbotapi.NewBotAPI(settings.Telegram)
	if fail != nil {
		log.Panic(fail)
	}
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 9
	channel := models.BotAPI.GetUpdatesChan(update)

	//db, err := sql.Open(settings.Database.Type, settings.Database.Arguments)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}

	dbpool, err := pgxpool.Connect(context.Background(), settings.Database.Arguments)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	iap := IAP.New()

	registerRepo := repository.NewRegisterRepo(*dbpool)

	registerUC := usecase2.NewRegisterUseCase(iap, registerRepo)

	for {
		message := <-channel
		mainMutex.Lock()
		chat, found := chats[message.FromChat().ID]
		if !found {
			chat = usecase.Chat{Id: message.FromChat().ID, Channel: make(chan tgbotapi.Update), Iap: iap, Register: registerUC}
			chats[message.FromChat().ID] = chat
			log.Println("Chat " + strconv.FormatInt(message.FromChat().ID, 10) + " created")
			go chat.Routine(chats, &mainMutex)
		}
		mainMutex.Unlock()
		chat.Channel <- message
	}
}
