package dao

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var sqlDB *sql.DB

func Start(conf config.Config) {
	c := conf.MysqlConf
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Pass, c.Host, c.Db)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("[db] database start fail: %s", err)
	}
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatalf("[db] database start fail: %s", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour * 4)
	sqlDB.SetConnMaxIdleTime(time.Hour * 4)

	query.SetDefault(db)
}
