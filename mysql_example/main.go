package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(constr string) {
	db, err := gorm.Open(mysql.New(mysql.Config{
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
	DB = db
}

type UserInfo struct {
	Uid  string
	Name string
}

func main() {
	user := UserInfo{}
	DB.First(&user, "uid = ?", 1)
	DB.Model(&user).Update("Name", "test")
	DB.Delete(&user, "uid = ?", 1)

	users := make([]UserInfo, 0)
	DB.Where("1=1").Scan(&users)
}
