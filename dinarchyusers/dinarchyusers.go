package dinarchyusers

import (
	"dinarchy/dcron"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

type DinarchyUsers struct {
	Dcrons map[int]dcron.Dcron
	Logger *zap.SugaredLogger
	Db     *gorm.DB
}

func (u *DinarchyUsers) MigrateDb() {
	u.Db.AutoMigrate(&dcron.Dcron{})
}

func (u *DinarchyUsers) GetDcron(user *tb.User) dcron.Dcron {
	if len(u.Dcrons) == 0 {
		u.Dcrons = make(map[int]dcron.Dcron)
	}
	if val, ok := u.Dcrons[user.ID]; ok {
		u.Logger.Debugw("GetDcron existing dcron", "key", user.ID, "val", val)
		return val
	}
	u.Dcrons[user.ID] = dcron.Dcron{Tguser: user, Cron: cron.New()}
	go u.Dcrons[user.ID].Cron.Start()
	u.Logger.Debugw("GetDcron new dcron", "key", user.ID, "val", u.Dcrons[user.ID])
	return u.Dcrons[user.ID]
}
