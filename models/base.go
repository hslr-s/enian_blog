package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beego/beego/v2/server/web"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type BaseModel struct {
	gorm.Model
}

// 分页的结构体
type PageLimitStruct struct {
	PageSize  int `gorm:"-"` //
	LimitSize int `gorm:"-"` //
}

// 计算分页
func calcPage(page_size, limit_size int) (offset, limit int) {
	offset = limit_size * (page_size - 1)
	limit = limit_size
	return
}

var Db *gorm.DB

func GetDb() (*gorm.DB, error) {
	dbDrive := web.AppConfig.DefaultString("database::drive", "sqlite")
	var db *gorm.DB
	var err error

	if dbDrive == "mysql" {
		fmt.Print("数据库驱动:", "MySQL", "\n")
		// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		// 	return "blog_" + defaultTableName
		// }

		// db, _ := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/enian_blog?charset=utf8mb4&parseTime=True&loc=Local")
		host := web.AppConfig.DefaultString("database::host", "127.0.0.1")
		port := web.AppConfig.DefaultString("database::port", "3306")
		db_name := web.AppConfig.DefaultString("database::database_name", "enian_blog")
		username := web.AppConfig.DefaultString("database::username", "root")
		password := web.AppConfig.DefaultString("database::password", "")
		dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db_name + "?charset=utf8mb4&parseTime=True&loc=Local"
		// fmt.Println("链接信息", dsn)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				// TablePrefix:   "blog_",
				SingularTable: true,
			},
			Logger:                                   GetLogger(),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		// db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
			&Article{},
			&User{},
			&Config{},
			&Anthology{},
			&Message{},
			&Tag{},
			&File{},
			&UserConfig{},
			&FriendLink{},
		)
		sqlDb, _ := db.DB()
		sqlDb.SetMaxIdleConns(10)             // SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDb.SetMaxOpenConns(100)            // SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDb.SetConnMaxLifetime(time.Minute) // SetConnMaxLifetime 设置了连接可复用的最大时间。

	} else {
		fmt.Println("数据库驱动:", "SQLite")
		db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				// TablePrefix:   "blog_",
				SingularTable: true,
			},
			Logger: GetLogger(),
			// DisableForeignKeyConstraintWhenMigrating: true,
		})

		db.AutoMigrate(
			&Article{},
			&User{},
			&Config{},
			&Anthology{},
			&Message{},
			&Tag{},
			&File{},
			&UserConfig{},
			&FriendLink{},
		)

	}

	Db = db
	return Db, err
}

// 日志
func GetLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Warn, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

}
