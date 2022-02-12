package telegrambot

import (
	"time"

	"github.com/rasulov-emirlan/pukbot/internal/puk"
	tele "gopkg.in/telebot.v3"
)

type Telegrambot struct {
	bot   *tele.Bot
	puk   puk.Service
	token string
}

func NewBot(botToken string, pukService puk.Service) (*Telegrambot, error) {
	pref := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	telebot := &Telegrambot{
		bot:   bot,
		token: botToken,
		puk:   pukService,
	}

	telebot.ConfigCommands()

	return telebot, nil
}

func (t *Telegrambot) Start() {
	t.bot.Start()
}

func (t *Telegrambot) Close() {
	t.bot.Stop()
}

func (t *Telegrambot) ConfigCommands() error {
	t.bot.SetWebhook(&tele.Webhook{})
	t.bot.Handle(tele.OnText, func(c tele.Context) error {
		return nil
	})
	t.bot.Handle("/upload", t.PukCreate())
	t.bot.Handle("/list", t.PukList)
	return nil
}
