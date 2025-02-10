package bot

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"reminder/internal/service"
	"strings"
	"sync"
	"time"
)

var userState = make(map[int64]string)
var mu sync.Mutex

func StartBot(token string) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	commands := []telebot.Command{
		{Text: "start", Description: "Запуск бота"},
		{Text: "add", Description: "После /add напишите заметку которую хотите добавить"},
		{Text: "notes", Description: "Вывод всех заметок"},
		{Text: "help", Description: "Информация о боте"},
	}

	err = bot.SetCommands(commands)
	if err != nil {
		log.Fatal("Ошибка вывода меню команд ", err)
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

	bot.Handle("/help", func(ctx telebot.Context) error {
		message := "📌 *Список доступных команд:*\n\n" +
			"🔹 `/start` - Запустить бота\n" +
			"🔹 `/add` - Добавить заметку\n" +
			"🔹 `/notes` - Показать все заметки\n" +
			"🔹 `/help` - Показать это сообщение\n" +
			"💡 Используйте кнопки ниже или введите команду вручную."
		return ctx.Send(message)
	})

	bot.Start()
}
