package bot

import "github.com/PaulSonOfLars/gotgbot/v2"

type Keyboards struct{}

func NewKeyboards() *Keyboards {
	return &Keyboards{}
}

func (k *Keyboards) MainMenu() [][]gotgbot.InlineKeyboardButton {
	return [][]gotgbot.InlineKeyboardButton{
		{
			{Text: "📚 Информация о курсе", CallbackData: "course_info"},
		},
		{
			{Text: "✅ Проверить подписку", CallbackData: "check_subscription"},
		},
	}
}

func (k *Keyboards) CourseMenu() [][]gotgbot.InlineKeyboardButton {
	return [][]gotgbot.InlineKeyboardButton{
		{
			{Text: "📝 Подробнее", CallbackData: "course_info"},
		},
		{
			{Text: "✅ Проверить подписку", CallbackData: "check_subscription"},
		},
	}
}

func (k *Keyboards) CourseDetailMenu() [][]gotgbot.InlineKeyboardButton {
	return [][]gotgbot.InlineKeyboardButton{
		{
			{Text: "✅ Проверить подписку", CallbackData: "check_subscription"},
		},
	}
}
