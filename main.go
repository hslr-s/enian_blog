package main

import (
	"enian_blog/lib/cmn"
	"enian_blog/models"
	_ "enian_blog/routers"
	"fmt"
	"os"

	"github.com/astaxie/beego"
)

func main() {
	// 创建表
	db, dbErr := models.GetDb()
	if dbErr != nil {
		fmt.Println("数据库错误：", dbErr)
		os.Exit(1)
	}

	//禁止表名自动加s
	// 创建表
	db.AutoMigrate(
		&models.Article{},
		&models.User{},
		&models.Config{},
		&models.Anthology{},
		&models.Tag{},
	)
	mUser := models.User{}
	err := db.Where("id=?", 1).First(&mUser).Error
	if err != nil {
		mUser.Mail = "admin@team.cc"
		mUser.Username = "admin"
		mUser.Password = cmn.PasswordEncryption("111111")
		mUser.Name = "超级用户"
		mUser.Status = 1
		mUser.Role = 1
		mUser.AddOne(mUser)
		fmt.Println("===============================================")
		fmt.Println("数据库初始化成功")
		fmt.Println("用户名：", "admin")
		fmt.Println("密码：", "111111")
		fmt.Println("邮箱（暂时可在数据库中修改）：", "admin@team.cc")
		fmt.Println("===============================================")
	}

	beego.Run()
}
