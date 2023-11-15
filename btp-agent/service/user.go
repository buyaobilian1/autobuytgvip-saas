package service

import (
	"errors"
	"log"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/query"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func FindUserByTgId(tgId int64) (*model.User, error) {
	var u = query.User
	dbUser, err := u.Where(u.TgID.Eq(tgId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("db query fail. %v\n", err)
		return nil, err
	}
	return dbUser, nil
}

func FindUserByBotId(botId int64) (*model.User, error) {
	var u = query.User
	dbUser, err := u.Where(u.BotID.Eq(botId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dbUser, nil
}

func FindOrCreateUserByTgCtx(ctx tele.Context) (*model.User, error) {
	var u = query.User
	user := &model.User{
		ParentID:   0,
		TgID:       ctx.Sender().ID,
		TgUsername: &ctx.Sender().Username,
		Balance:    0,
		Brokerage:  0,
		TronAddr:   nil,
		CreatedAt:  time.Time{},
		BotStatus:  2,
	}
	dbUser, err := u.Attrs(field.Attrs(user)).Where(u.TgID.Eq(ctx.Sender().ID)).FirstOrCreate()
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}
