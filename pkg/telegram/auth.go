package telegram

import (
	"context"
	"fmt"

	"github.com/Oden333/TG_bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authlink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf(replyStartTemplate, authlink))

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getAccessToken(chatId int64) (string, error) {
	return b.tokenRepository.Get(chatId, repository.AccessTokens)
}

func (b *Bot) generateAuthorizationLink(chatId int64) (string, error) {
	redirectUrl := b.generateRedirectURL(chatId)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatId, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}
	return b.pocketClient.GetAuthorizationURL(requestToken, redirectUrl)
}

func (b *Bot) generateRedirectURL(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
