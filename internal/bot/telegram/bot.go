package telegram

import (
	"fmt"
	"moscow-events-telegramauth/internal/bot/consts"
	"moscow-events-telegramauth/internal/config"
	"moscow-events-telegramauth/internal/storage"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go.uber.org/zap"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	redisStorage storage.RedisStorage
	logger       *zap.SugaredLogger
	timeOut      int
}

func InitBot(
	cfg *config.Config,
	redisStorage storage.RedisStorage,
	logger *zap.SugaredLogger,
	debug bool) (*Bot, error) {
	const op = "bot.InitBot"

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBot.Token)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	tgBot := Bot{
		bot:          bot,
		redisStorage: redisStorage,
		logger:       logger,
		timeOut:      cfg.TelegramBot.TimeOut,
	}
	bot.Debug = debug
	logger.Infof("Authorized on account %s", bot.Self.UserName)
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
		b.logger.Debug(upd.CallbackQuery)
	}
	if mes != nil {
		b.logger.Debug(mes.Text)
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
