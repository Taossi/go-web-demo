package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

/*
  利用gorm.Model
  包含 ID, CreatedAt, UpdatedAt, DeletedAt
*/
type User struct {
	gorm.Model
	//'gorm:"type:varchar(20);not null"'
	Name     string
	Phone    string
	Password string
}

// 需要先在本机用Navicat创建一个名为go_db的数据库
// 注意修改用户密码
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
	//这里 gorm.Open()函数与之前版本的不一样，大家注意查看官方最新gorm版本的用法
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error to Db connection, err: " + err.Error())
	}
	//这个是gorm自动创建数据表的函数。它会自动在数据库中创建一个名为users的数据表
	_ = db.AutoMigrate(&User{})
	return db
}

func main() {
	db := InitDB()

	//利用Gin框架的web写法，来源于gin官网
	r := gin.Default()
	r.POST("/api/auth/register", func(c *gin.Context) {
		//获取参数
		name := c.PostForm("name")
		phone := c.PostForm("phone")
		password := c.PostForm("password")
		//数据验证
		if len(phone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "手机号格式不正确",
			})
			return
		}

		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码不能少于6位",
			})
			return
		}
		if len(name) == 0 {
			name = RandomString(10)
			return
		}
		log.Print(name, phone, password)

		//判断手机号是否存在
		if isPhoneExist(db, phone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "手机号已被注册，请重试",
			})
			return
		}

		//创建新用户 插入数据
		newUser := User{
			Name:     name,
			Phone:    phone,
			Password: password,
		}
		db.Create(&newUser)

		//返回结果
		c.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	_ = r.Run() // listen and serve on 0.0.0.0:8080
	panic(r.Run())
}

/*
	查询数据库中是否存在当前电话
*/
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user User
	db.Where("phone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

//随机产生英文字符，可设定长度
func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
