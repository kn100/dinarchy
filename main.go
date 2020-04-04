package main

import (
	"dinarchy/models"
	"dinarchy/services"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

	db, err := gorm.Open("sqlite3", "/database/dinarchy.sqlite3") //Add to env file
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&models.Job{})

	// You can get an API key from the botfather...
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TGBOT_KEY"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		sugar.Fatal(err)
	}

	sugar.Debug("Created Telegram Service")

	cs := services.CronService{TS: b}
	cs.Init()

	js := services.JobService{DB: db, CS: cs}
	js.LoadJobs()

	sugar.Debug("Created Job Service")

	b.Handle("/start", func(m *tb.Message) {
		handleHelp(b, m.Sender.ID, m.Text)
	})

	b.Handle("/help", func(m *tb.Message) {
		handleHelp(b, m.Sender.ID, m.Text)
		sugar.Debugw("/help", "args", m.Payload, "sender", m.Sender)
	})

	b.Handle("/create", func(m *tb.Message) {
		handleCreate(b, js, m.Sender.ID, m.Payload)
		sugar.Debugw("/create", "args", m.Payload, "sender", m.Sender)
	})

	b.Handle("/show", func(m *tb.Message) {
		handleShow(b, js, m.Sender.ID, m.Payload)
		sugar.Debugw("/show", "args", m.Payload, "sender", m.Sender)
	})

	b.Handle("/delete", func(m *tb.Message) {
		handleDelete(b, js, m.Sender.ID, m.Payload)
		sugar.Debugw("/delete", "args", m.Payload, "sender", m.Sender)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		handleOther(b, m.Sender.ID, m.Payload)
		sugar.Debugw("unknown", "args", m.Payload, "sender", m.Sender)
	})

	sugar.Infof("Bot awaiting commands")
	b.Start()
}

func handleCreate(tgs *tb.Bot, js services.JobService, tgid int, text string) {
	tgid_str := strconv.Itoa(tgid)
	cmd := strings.TrimPrefix(text, "/create ")

	argsplit := split(cmd)
	if len(argsplit) < 3 {
		handleOther(tgs, tgid, text)
		return
	}
	err := js.AddJob(tgid_str, argsplit[0], argsplit[1], argsplit[2])
	if err != nil {
		tgs.Send(&services.Recipient{TGID: tgid}, "Could not create cron: "+err.Error(), &tb.SendOptions{ParseMode: tb.ModeMarkdown})
		return
	}
	tgs.Send(&services.Recipient{TGID: tgid}, fmt.Sprintf("Created job with name %s", argsplit[1]), &tb.SendOptions{ParseMode: tb.ModeMarkdown})

}

func handleShow(tgs *tb.Bot, js services.JobService, tgid int, _ string) {
	tgid_str := strconv.Itoa(tgid)

	jobs := js.GetJobs(tgid_str)
	jobstr := fmt.Sprintf("Your %d jobs: \n", len(jobs))
	for _, j := range jobs {
		jobstr += fmt.Sprintf("cron: `%s`, name: `%s`\n", j.CronString, j.Name)
	}
	fmt.Println("jobstr:", string(jobstr))
	tgs.Send(&services.Recipient{TGID: tgid}, jobstr, &tb.SendOptions{ParseMode: tb.ModeMarkdown})

}

func handleDelete(tgs *tb.Bot, js services.JobService, tgid int, text string) {
	//Todo: pass around the recipient object, rather than this fucking tgid
	tgid_str := strconv.Itoa(tgid)
	cmd := strings.TrimPrefix(text, "/delete ")

	if err := js.RemoveJob(tgid_str, cmd); err != nil {
		tgs.Send(&services.Recipient{TGID: tgid}, "Could not remove cron: "+err.Error(), &tb.SendOptions{ParseMode: tb.ModeMarkdown})

		return
	}
	tgs.Send(&services.Recipient{TGID: tgid}, fmt.Sprintf("Removed cron %s", cmd))
}

func handleHelp(tgs *tb.Bot, tgid int, _ string) {
	s := "Dinarchy is a bot for scheduling reminders using Cron syntax.\n\n` /create 0 9 * * *,milk-check,Check the milk` would schedule a reminder to check the milk at 09:00 every morning, called milk-check for example\n\n`/delete milk-check` will delete the cron job with the name milk-check."
	s2 := "/show shows your commands"
	tgs.Send(&services.Recipient{TGID: tgid}, s, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	tgs.Send(&services.Recipient{TGID: tgid}, s2, &tb.SendOptions{ParseMode: tb.ModeMarkdown})

}

func handleOther(tgs *tb.Bot, tgid int, _ string) {
	tgs.Send(&services.Recipient{TGID: tgid}, "Unknown command. try /help", &tb.SendOptions{ParseMode: tb.ModeMarkdown})
}

func split(cmd string) []string {
	splitFunc := func(c rune) bool {
		return c == ','
	}
	return strings.FieldsFunc(cmd, splitFunc)
}
