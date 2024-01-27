package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/UdinSemen/moscow-events-telegramauth/internal/bot/consts"
	pg_storage "github.com/UdinSemen/moscow-events-telegramauth/internal/storage/pg-storage"
	"github.com/UdinSemen/moscow-events-telegramauth/internal/storage/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

const (
	timeCodeLen = 6
)

func (b *Bot) HandleText(upd tgbotapi.Update) {
	text := upd.Message.Text
	switch len(text) {
	case timeCodeLen:
		b.handleTimeCode(upd)

	default:
		b.handleInvalid(upd)
	}
	if len(text) == timeCodeLen {

	}
}

func (b *Bot) handleTimeCode(upd tgbotapi.Update) {
	const op = "telegram.handleTimeCode"
	timeCode := upd.Message.Text
	chatID := upd.Message.Chat.ID

	for _, l := range timeCode {
		if l < '0' || l > '9' {
			send := tgbotapi.NewMessage(chatID, consts.TextMustHaveOnlyLetters)
			if _, err := b.bot.Send(send); err != nil {
				zap.S().Errorf("%s:%s", op, err)
			}
			return
		}
	}

	firstName := upd.Message.From.FirstName
	lastName := upd.Message.From.LastName
	userID := upd.Message.From.ID

	flagHaveUser := false
	if err := b.pgStorage.AddUser(firstName, lastName, "", userID); err != nil {
		if errors.Is(err, pg_storage.ErrPgUniqueConstr) {
			flagHaveUser = true
		} else {
			zap.S().Errorf("%s:%s", op, err)
			if err := b.sentSmtWrongWithReq(chatID); err != nil {
				zap.S().Errorf("%s:%s", op, err)
			}
			return
		}
	}

	if err := b.redisStorage.ConfirmSession(context.Background(), timeCode,
		strconv.FormatInt(upd.Message.Chat.ID, 10)); err != nil {
		zap.S().Error(fmt.Errorf("%s:%w", op, err))
		if errors.Is(err, redis.ErrEmptyJSONValue) {
			send := tgbotapi.NewMessage(upd.Message.Chat.ID, consts.TextErrEmpty)
			if _, err := b.bot.Send(send); err != nil {
				zap.S().Errorf("%s:%s", op, err)
			}
			return
		}
		if errors.Is(err, redis.ErrSessionAlreadyConfirmed) {
			send := tgbotapi.NewMessage(upd.Message.Chat.ID, consts.TextErrAlreadyConfirmed)
			if _, err := b.bot.Send(send); err != nil {
				zap.S().Errorf("%s:%s", op, err)
			}
			return
		}

		if err := b.sentSmtWrongWithReq(chatID); err != nil {
			zap.S().Error(fmt.Errorf("%s:%w", op, err))
			return
		}
		return
	}

	textMsg := consts.TextConfirmed
	if flagHaveUser {
		textMsg = consts.TextAlreadyHaveUser
	}

	send := tgbotapi.NewMessage(chatID, textMsg)
	if _, err := b.bot.Send(send); err != nil {
		zap.S().Errorf("%s:%s", op, err)
		return
	}

}

func (b *Bot) handleInvalid(upd tgbotapi.Update) {
	const op = "telegram.handleTimeCode"

	send := tgbotapi.NewMessage(upd.Message.Chat.ID, fmt.Sprintf(consts.TextCodeMustHaveSixLen, timeCodeLen))
	if _, err := b.bot.Send(send); err != nil {
		zap.S().Errorf("%s:%s", op, err)
	}
}
