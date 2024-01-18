package telegram

import (
	"fmt"

	"github.com/UdinSemen/moscow-events-telegramauth/internal/bot/consts"
	"github.com/UdinSemen/moscow-events-telegramauth/internal/config"
	"github.com/UdinSemen/moscow-events-telegramauth/internal/storage"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go.uber.org/zap"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	redisStorage storage.RedisStorage
	pgStorage    storage.PgStorage
	timeOut      int
}

func InitBot(
	cfg *config.Config,
	redisStorage storage.RedisStorage,
	pgStorage storage.PgStorage,
	debug bool) (*Bot, error) {
	const op = "bot.InitBot"

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBot.Token)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	tgBot := Bot{
		bot:          bot,
		redisStorage: redisStorage,
		pgStorage:    pgStorage,
		timeOut:      cfg.TelegramBot.TimeOut,
	}
	bot.Debug = debug
	zap.S().Infof("Authorized on account %s", bot.Self.UserName)
	return &tgBot, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.timeOut

	updates := b.bot.GetUpdatesChan(u)
	for upd := range updates {
		b.handleUpdate(upd)
	}
}

func (b *Bot) handleUpdate(upd tgbotapi.Update) {

	mes := upd.Message
	if upd.CallbackQuery != nil {
		zap.S().Debug(upd.CallbackQuery)
	}
	if mes != nil {
		zap.S().Debug(mes.Text)
		if mes.IsCommand() {
			b.HandleCommand(upd)
		}
		if !mes.IsCommand() {
			b.HandleText(upd)
		}
	}
}

func (b *Bot) sentSmtWrongWithReq(chatID int64) error {
	send := tgbotapi.NewMessage(chatID, consts.TextSmtWrongWithRequest)
	if _, err := b.bot.Send(send); err != nil {
		return err
	}
	return nil
}
