package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnauthurized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save url")
)

func (b *Bot) handleError(ChatID int64, err error) {
	msg := tgbotapi.NewMessage(ChatID, "Произошла неизвестная ошибка")

	switch err {
	case errInvalidURL:
		msg.Text = "Неверная ссылка"
		b.bot.Send(msg)
	case errUnauthurized:
		msg.Text = "Неверная авторизация(Попробуй /start)"
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = "Сохранение не удалось(Попробуй ещё раз)"
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}

}
