package service

import (
	"errors"
	"log"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
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
		log.Printf("db query fail. %v\n", err)
		return nil, err
	}
	return dbUser, nil
}

func FindOrCreateUserByTgCtx(ctx tele.Context) (*model.User, error) {
	var u = query.User
	agentUser, err := u.Where(u.BotID.Eq(ctx.Bot().Me.ID)).First()
	if err != nil {
		log.Printf("[db] query fail. %v\n", err)
		return nil, err
	}
	user := &model.User{
		ParentID:   agentUser.ID,
		TgID:       ctx.Sender().ID,
		TgUsername: &ctx.Sender().Username,
		Balance:    0,
		Brokerage:  0,
		TronAddr:   nil,
		CreatedAt:  time.Time{},
	}
	dbUser, err := u.Attrs(field.Attrs(user)).Where(u.TgID.Eq(ctx.Sender().ID)).FirstOrCreate()
	if err != nil {
		log.Printf("[db] query fail. %v\n", err)
		return nil, err
	}
	return dbUser, nil
}
