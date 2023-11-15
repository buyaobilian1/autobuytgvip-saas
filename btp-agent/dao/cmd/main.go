package main

import (
	"github.com/buyaobilian1/autobuytgvip-saas/btp-agent/dao/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const mysqlDns = "btp_saas:4SpcrSai4cHZxbRi@(127.0.0.1:3306)/btp_saas?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../query",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		FieldNullable: true,
	})

	gormdb, _ := gorm.Open(mysql.Open(mysqlDns))
	g.UseDB(gormdb) // reuse your gorm db
	g.ApplyBasic(model.Order{}, model.User{}, model.Recharge{}, model.Param{}, model.Withdraw{})
	g.Execute()
}
