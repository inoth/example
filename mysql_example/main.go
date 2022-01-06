package main

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDatabase(constr string) {
	conn, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       constr,
		DefaultStringSize:         1024, // string 类型字段的默认长度
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	sqldb, err := conn.DB()
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	sqldb.SetConnMaxIdleTime(10)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetConnMaxLifetime(time.Second * 60)
	db = conn
}
func GetDb() *gorm.DB {
	return db
}

type UserInfo struct {
	Uid  string
	Name string
}

func main() {
	InitDatabase("user:passwd@(localhost:3306)/user?charset=utf8&parseTime=True&loc=Local")

	mysqldb := GetDb()
	user := UserInfo{}
	mysqldb.First(&user, "uid = ?", 1)
	mysqldb.Model(&user).Update("Name", "test")
	mysqldb.Delete(&user, "uid = ?", 1)

	users := make([]UserInfo, 0)
	mysqldb.Where("1=1").Scan(&users)
}
