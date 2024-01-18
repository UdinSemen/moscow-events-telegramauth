package telegram

import (
	"github.com/UdinSemen/moscow-events-telegramauth/internal/bot/consts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleCommand(upd tgbotapi.Update) {
	switch upd.Message.Command() {
	case "start":
		b.startCommand(upd)
	}
}

func (b *Bot) startCommand(upd tgbotapi.Update) {
	chatID := upd.Message.Chat.ID
	send := tgbotapi.NewMessage(chatID, consts.HelloMsg)
	if _, err := b.bot.Send(send); err != nil {
		b.logger.Error(err)
	}
}
