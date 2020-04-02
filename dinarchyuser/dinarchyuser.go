package dinarchyuser

import (
	"dinarchy/dcron"

	tb "gopkg.in/tucnak/telebot.v2"
)

type DinarchyUser struct {
	Tguser *tb.User
	Dcron  *dcron.Dcron
}
