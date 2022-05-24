package main

import (
	"enian_blog/lib/cmn"
	"enian_blog/lib/initialize"
	"enian_blog/models"
	_ "enian_blog/routers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	// beego.BConfig.WebConfig.Session.SessionOn = true
	// beego.BConfig.WebConfig.Session.SessionProvider = "file"
	// beego.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
}

func main() {
	initialize.ConfigCreate()
	initialize.PrintlnLogoAndVersion()
	// 创建表
	db, dbErr := models.GetDb()
	if dbErr != nil {
		cmn.FatalError("数据库错误," + dbErr.Error())
	}
	// db.AutoMigrate(
	// 	&models.Article{},
	// 	&models.User{},
	// 	&models.Config{},
	// 	&models.Anthology{},
	// 	&models.Message{},
	// 	&models.Tag{},
	// 	&models.File{},
	// )

	mUser := models.User{}
	err := db.First(&mUser).Error
	if err != nil {
		initialize.CreateData()
	}
	web.AddFuncMap("TimeToRelativeTime", cmn.TimeToRelativeTime)
	web.AddFuncMap("TimeStrToRelativeTime", cmn.TimeStrToRelativeTime)
	web.Run()
}
