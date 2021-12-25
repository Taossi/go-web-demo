package controller

import (
	"gin-gorm/common"
	"gin-gorm/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

/**
用户注册逻辑函数
*/
func UserRegister(c *gin.Context) {
	db := common.InitDB()
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
	newUser := model.User{
		Name:     name,
		Phone:    phone,
		Password: password,
	}
	db.Create(&newUser)

	//返回结果
	c.JSON(200, gin.H{
		"message": "注册成功",
	})
	return
}

/**
查询数据库中是否存在当前电话
*/
func isPhoneExist(db *gorm.DB, phone string) bool {
	var user model.User
	db.Where("phone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

/**
随机产生指定长度的英文字符作为用户名
*/
func RandomString(n int) string {
	letters := []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
