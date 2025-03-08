package bot

import (
	"cursebot/internal/repository"
	"fmt"
	"log"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type Bot struct {
	bot       *gotgbot.Bot
	repo      repository.Repository
	channelID string
	keyboards *Keyboards
}

func NewBot(bot *gotgbot.Bot, repo repository.Repository, channelID string) *Bot {
	return &Bot{
		bot:       bot,
		repo:      repo,
		channelID: channelID,
		keyboards: NewKeyboards(),
	}
}

func (b *Bot) Start() error {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Printf("Error handling update: %v", err)
			return ext.DispatcherActionNoop
		},
	})

	// Регистрация обработчиков команд
	dispatcher.AddHandler(handlers.NewCommand("start", b.handleStart))
	dispatcher.AddHandler(handlers.NewCommand("help", b.handleHelp))
	dispatcher.AddHandler(handlers.NewCommand("course", b.handleCourse))

	// Регистрация обработчиков для callback'ов
	dispatcher.AddHandler(handlers.NewCallback(nil, b.handleCallbackQuery))

	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		ErrorLog: nil,
	})

	err := updater.StartPolling(b.bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 60,
			AllowedUpdates: []string{
				"message",
				"callback_query",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to start polling: %w", err)
	}

	log.Printf("Bot %s started", b.bot.User.Username)
	updater.Idle()
	return nil
}

func (b *Bot) handleStart(bot *gotgbot.Bot, ctx *ext.Context) error {
	channelID, _ := strconv.ParseInt(b.channelID, 10, 64)
	chatMember, err := bot.GetChatMember(channelID, ctx.EffectiveUser.Id, nil)

	var text string
	if err != nil {
		log.Printf("Error getting chat member: %v", err)
		text = "Произошла ошибка при проверке подписки. Пожалуйста, попробуйте позже."
	} else {
		status := chatMember.GetStatus()
		log.Printf("User %d subscription status: %s", ctx.EffectiveUser.Id, status)

		switch status {
		case "member", "administrator", "creator":
			text = fmt.Sprintf("Привет, %s! 👋\n\nВы успешно подписаны на канал. Теперь вам доступны следующие функции:", ctx.EffectiveUser.FirstName)
			_, err = ctx.EffectiveMessage.Reply(bot, text, &gotgbot.SendMessageOpts{
				ParseMode: "Markdown",
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
						{
							Text:         "Информация о курсе",
							CallbackData: "course_info",
						},
						{
							Text:         "Приобрести курс",
							CallbackData: "purchase_course",
						},
					}},
				},
			})
			return err
		default:
			text = "Вы не подписаны на канал. Пожалуйста, подпишитесь на наш канал, чтобы получить доступ к курсу: [наш канал](https://t.me/vibecodinghub1)"
		}
	}

	_, err = ctx.EffectiveMessage.Reply(bot, text, &gotgbot.SendMessageOpts{
		ParseMode: "Markdown",
	})
	return err
}

func (b *Bot) handleHelp(bot *gotgbot.Bot, ctx *ext.Context) error {
	text := `Доступные команды:
/start - Начать работу с ботом
/help - Показать это сообщение`

	_, err := ctx.EffectiveMessage.Reply(bot, text, nil)
	return err
}

func (b *Bot) handleCallbackQuery(bot *gotgbot.Bot, ctx *ext.Context) error {
	log.Printf("Processing callback query with data: %s", ctx.CallbackQuery.Data)

	switch ctx.CallbackQuery.Data {
	case "check_subscription":
		channelID, _ := strconv.ParseInt(b.channelID, 10, 64)
		chatMember, err := bot.GetChatMember(channelID, ctx.CallbackQuery.From.Id, nil)

		var text string
		if err != nil || chatMember.GetStatus() == "left" || chatMember.GetStatus() == "kicked" {
			text = "Вы не подписаны на канал. Пожалуйста, подпишитесь для доступа к курсу."
		} else {
			text = fmt.Sprintf("Привет, %s! 👋\n\nВы успешно подписаны на канал. Теперь вам доступны следующие функции:", ctx.EffectiveUser.FirstName)

			_, _, err = ctx.CallbackQuery.Message.EditText(bot, text, &gotgbot.EditMessageTextOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
						{
							Text:         "Информация о курсе",
							CallbackData: "course_info",
						},
						{
							Text:         "Приобрести курс",
							CallbackData: "purchase_course",
						},
					}},
				},
			})
			return err
		}

		_, err = ctx.CallbackQuery.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: text,
		})
		return err

	case "course_info":
		courses, err := b.repo.ListCourses()
		if err != nil {
			log.Printf("Error getting courses: %v", err)
			return err
		}

		if len(courses) == 0 {
			log.Printf("No courses available")
			return fmt.Errorf("no courses available")
		}

		course := courses[0]
		text := fmt.Sprintf("*Подробная информация о курсе*\n\n"+
			"*%s*\n\n%s\n\n"+
			"*Что вы получите:*\n"+
			"✅ Доступ к материалам курса\n"+
			"✅ Практические задания\n"+
			"✅ Поддержку кураторов\n"+
			"✅ Сертификат о прохождении\n\n"+
			"Цена: %.2f руб.", course.Title, course.Description, course.Price)

		_, err = bot.SendMessage(ctx.EffectiveChat.Id, text, &gotgbot.SendMessageOpts{
			ParseMode: "Markdown",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: b.keyboards.CourseDetailMenu(),
			},
		})
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}

		_, err = ctx.CallbackQuery.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
			Text: "Информация о курсе отправлена",
		})
		return err

	default:
		log.Printf("Unknown callback data: %s", ctx.CallbackQuery.Data)
		return nil
	}
}

func (b *Bot) handleCourse(bot *gotgbot.Bot, ctx *ext.Context) error {
	courses, err := b.repo.ListCourses()
	if err != nil {
		log.Printf("Error getting courses: %v", err)
		_, err = ctx.EffectiveMessage.Reply(bot, "Произошла ошибка при получении информации о курсе", nil)
		return err
	}

	if len(courses) == 0 {
		_, err = ctx.EffectiveMessage.Reply(bot, "В данный момент нет доступных курсов", nil)
		return err
	}

	course := courses[0] // Для примера берем первый курс
	text := fmt.Sprintf("*%s*\n\n%s\n\nЦена: %.2f руб.",
		course.Title, course.Description, course.Price)

	_, err = ctx.EffectiveMessage.Reply(bot, text, &gotgbot.SendMessageOpts{
		ParseMode: "Markdown",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: b.keyboards.CourseMenu(),
		},
	})
	return err
}
