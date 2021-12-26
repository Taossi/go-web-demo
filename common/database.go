package common

import (
	"fmt"
	"gin-gorm/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

/**
 * @Description: MySQL数据库初始化 需要修改相应的用户与密码
 * @return *gorm.DB
 */
func InitDB() *gorm.DB {
	host := "localhost"
	port := "3306"
	database := "go_db"
	username := "root"
	password := "908109ywshabeer"
	charset := "utf8"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error to Db connection, err: " + err.Error())
	}
	//这个是gorm自动创建数据表的函数。它会自动在数据库中创建一个名为users的数据表
	_ = db.AutoMigrate(&model.User{})
	return db
}

func getDB() *gorm.DB {
	return DB
}
