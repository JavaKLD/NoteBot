package bot

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"reminder/internal/service"
	"strings"
	"time"
)

func StartBot(token string) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(ctx telebot.Context) error {
		userId := ctx.Sender().ID
		userName := ctx.Sender().FirstName
		err := service.AddUser(userId, userName)
		if err != nil {
			return ctx.Send( "Ошибка добавления пользователя")
		}
		return ctx.Send("Добро пожаловать в бота для создания заметок, ")
	})
	
	bot.Handle("/add", func(ctx telebot.Context) error {
		arg := ctx.Args()
		if len(arg) == 0 {
			return ctx.Send("Введите текст заметки после командв /add")
		}

		content := strings.Join(arg, " ")
		err := service.CreateNote(ctx.Sender().ID, content)
		if err != nil {
			return ctx.Send("Ошибка при сохранении заметок.")
		}

		return ctx.Send("Заметка сохранена")
	})

	bot.Handle("/notes", func(ctx telebot.Context) error {
		notes, err := service.GetNotes(ctx.Sender().ID)
		if err != nil {
			return ctx.Send("Ошибка при получении заметок")
		}

		if len(notes) == 0 {
			return ctx.Send( "У вас нет заметок")
		}

		response := "Ваши заметки:\n"
		for i, note := range notes {
			response += fmt.Sprintf("%d. %s\n", i + 1, note)
		}
		return ctx.Send(response)
	})

	bot.Start()
}
