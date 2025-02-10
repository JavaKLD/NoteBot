package bot

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"reminder/internal/service"
	"strconv"
	"strings"
	"sync"
	"time"
)

var userState = make(map[int64]string)
var userNoteId = make(map[int64]int64)
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
		{Text: "add", Description: "Добавление заметки"},
		{Text: "notes", Description: "Вывод всех заметок"},
		{Text: "/reduct", Description: "Редактирование заметки"},
		{Text: "delete", Description: "Удаление выбранной заметки"},
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
		return ctx.Send("Добро пожаловать в бота для создания заметок. ")
	})
	
	bot.Handle("/add", func(ctx telebot.Context) error {
		mu.Lock()
		userState[ctx.Sender().ID] = "waiting_for_note"
		mu.Unlock()

		return ctx.Send("✏️ Введите текст заметки:")
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

	bot.Handle("/reduct", func(ctx telebot.Context) error {
		mu.Lock()
		userState[ctx.Sender().ID] = "waiting_for_reduct_id"
		mu.Unlock()

		notes, err := service.GetNotes(ctx.Sender().ID)
		if err != nil {
			return ctx.Send("Ошибка при получении заметок.")
		}

		if len(notes) == 0 {
			return ctx.Send("У вас нет заметок для обновления.")
		}
		response := "Ваши заметки:\n"
		for i, note := range notes {
			response += fmt.Sprintf("%d. %s\n", i+1, note)
		}
		response += "\nВведите номер заметки для обновления."

		return ctx.Send(response)
	})
	
	bot.Handle("/delete", func(ctx telebot.Context) error {
		mu.Lock()
		userState[ctx.Sender().ID] = "waiting_for_delete"
		mu.Unlock()

		notes, err := service.GetNotes(ctx.Sender().ID)
		if err != nil {
			return ctx.Send("Ошибка при получении заметок.")
		}

		if len(notes) == 0 {
			return ctx.Send("У вас нет заметок для удаления.")
		}
		response := "Ваши заметки:\n"
		for i, note := range notes {
			response += fmt.Sprintf("%d. %s\n", i+1, note)
		}
		response += "\nВведите номер заметки для удаления."

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

	bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		mu.Lock()
		state, exists := userState[ctx.Sender().ID]
		mu.Unlock()

		if exists && state == "waiting_for_note" {
			mu.Lock()
			delete(userState, ctx.Sender().ID)
			mu.Unlock()

			content := strings.TrimSpace(ctx.Text())
			err := service.CreateNote(ctx.Sender().ID, content)
			if err != nil {
				log.Fatal("❌ Ошибка при сохранении заметки.", err)
			}
			return ctx.Send("✅ Ваша заметка сохранена.")
		}

		if exists && state == "waiting_for_delete" {
			mu.Lock()
			delete(userState, ctx.Sender().ID)
			mu.Unlock()

			noteIndex, err := strconv.ParseInt(strings.TrimSpace(ctx.Text()), 10, 64)
			if err != nil || noteIndex < 1 {
				return ctx.Send("❌ Введите правильный номер заметки для удаления.")
			}

			err = service.DeleteNote(ctx.Sender().ID, noteIndex)
			if err != nil {
				return ctx.Send("❌ Ошибка при удалении заметки.")
			}

			return ctx.Send("✅ Заметка удалена.")
		}

		if exists && state == "waiting_for_reduct_id" {
			noteIndex, err := strconv.Atoi(strings.TrimSpace(ctx.Text()))
			if err != nil || noteIndex < 1 {
				return ctx.Send("❌ Введите правильный номер заметки.")
			}

			mu.Lock()
			userState[ctx.Sender().ID] = "waiting_for_reduct_text"
			userNoteId[ctx.Sender().ID] = int64(noteIndex)
			mu.Unlock()

			return ctx.Send("✏️ Введите новый текст заметки:")
		}

		if exists && state == "waiting_for_reduct_text" {
			mu.Lock()
			noteId := userNoteId[ctx.Sender().ID]
			delete(userNoteId, ctx.Sender().ID)
			delete(userState, ctx.Sender().ID)
			mu.Unlock()

			content := strings.TrimSpace(ctx.Text())
			err := service.RedNote(ctx.Sender().ID, noteId, content)
			if err != nil {
				return ctx.Send("❌ Ошибка при обновлении заметки.")
			}
			return ctx.Send("✅ Заметка успешно обновлена.")
		}
		return nil
	})

	bot.Start()
}
