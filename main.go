package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"dinarchy/dcron"
	"dinarchy/dinarchyusers"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	err := godotenv.Load()
	if err != nil {
		sugar.Fatal("error loading .env file")
	}

	// You can get an API key from the botfather...
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TGBOT_KEY"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		sugar.Fatal(err)
	}
	sugar.Debug("Started Telegram Bot")

	dinarchyusers := dinarchyusers.DinarchyUsers{Logger: sugar}
	sugar.Debug("Created user store")

	b.Handle("/start", func(m *tb.Message) {
		handleHelp(b, dinarchyusers.GetDcron(m.Sender), m.Text)
	})

	b.Handle("/help", func(m *tb.Message) {
		handleHelp(b, dinarchyusers.GetDcron(m.Sender), m.Payload)
		sugar.Debugw("/help", "args", m.Payload, "sender", m.Sender)
	})

	b.Handle("/create", func(m *tb.Message) {
		handleCreate(b, dinarchyusers.GetDcron(m.Sender), m.Payload)
		sugar.Debugw("/create", "args", m.Payload, "sender", m.Sender)
	})

	b.Handle("/delete", func(m *tb.Message) {
		handleDelete(b, dinarchyusers.GetDcron(m.Sender), m.Payload)
		sugar.Debugw("/delete", "args", m.Payload, "sender", m.Sender)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		handleOther(b, dinarchyusers.GetDcron(m.Sender), m.Payload)
		sugar.Debugw("unknown", "args", m.Payload, "sender", m.Sender)
	})

	sugar.Infof("Bot awaiting commands")
	b.Start()
}

func handleCreate(b *tb.Bot, dc dcron.Dcron, c string) {

	cmd := strings.TrimPrefix(c, "/create ")

	argsplit := split(cmd)
	if len(argsplit) < 3 {
		handleOther(b, dc, c)
		return
	}
	_, err := dc.AddJob(argsplit[0], argsplit[1], argsplit[2], b.Send)
	if err != nil {
		b.Send(dc.Tguser, "Could not create cron: "+err.Error())
		return
	}
	b.Send(dc.Tguser, fmt.Sprintf("Created job with name %s", argsplit[1]))
}

func handleDelete(b *tb.Bot, dc dcron.Dcron, c string) {
	cmd := strings.TrimPrefix(c, "/delete ")

	if err := dc.RemoveJob(cmd); err != nil {
		b.Send(dc.Tguser, "Could not remove cron: "+err.Error())
		return
	}
	b.Send(dc.Tguser, fmt.Sprintf("Removed cron %s", cmd))
}

func handleHelp(b *tb.Bot, dc dcron.Dcron, _ string) {
	s := "Dinarchy is a bot for scheduling reminders using Cron syntax.\n\n` /create 0 9 * * *,milk-check,Check the milk` would schedule a reminder to check the milk at 09:00 every morning, called milk-check for example\n\n`/delete milk-check` will delete the cron job with the name milk-check."
	b.Send(dc.Tguser, s, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
}

func handleOther(b *tb.Bot, dc dcron.Dcron, _ string) {
	b.Send(dc.Tguser, "Unknown command. try /help")
}

func split(cmd string) []string {
	splitFunc := func(c rune) bool {
		return c == ','
	}
	return strings.FieldsFunc(cmd, splitFunc)
}
