package controller

import (
	"gin-gorm/common"
	"gin-gorm/dto"
	"gin-gorm/model"
	"gin-gorm/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

/**
 * @Description: 用户注册逻辑函数 POST请求
 * @param c
 */
func UserRegister(c *gin.Context) {
	db := common.InitDB()
	//获取参数
	name := c.PostForm("name")
	phone := c.PostForm("phone")
	password := c.PostForm("password")

	//数据验证
	if len(phone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号格式不正确")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	if len(name) == 0 {
		name = RandomString(10)
	}
	log.Print(name, phone, password)

	//判断手机号是否存在
	if isPhoneExist(db, phone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号已被注册，请重试")
		return
	}

	//密码无法明文保存，需要加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "加密失败")
		return
	}
	log.Print(name, phone, hashedPassword)

	//创建新用户 插入数据
	newUser := model.User{
		Name:     name,
		Phone:    phone,
		Password: string(hashedPassword),
	}
	db.Create(&newUser)

	//返回结果
	response.Success(c, nil, "注册成功")
	return
}

/**
 * @Description: 用户登陆函数 通过phone, password登陆 POST请求
 * @param c
 */
func Login(c *gin.Context) {
	db := common.InitDB()
	//获取参数
	phone := c.PostForm("phone")
	password := c.PostForm("password")

	//数据验证
	if len(phone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号格式不正确")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	//判断手机号是否存在数据库
	var user model.User
	db.Where("phone = ?", phone).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该用户不存在")
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(c, nil, "密码错误")
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Fail(c, nil, "token错误")
		return
	}

	//返回结果
	response.Success(c, gin.H{"token": token}, "登陆成功")
	return
}

/**
 * @Description: 从context上下文获取用户信息(此时已通过token验证)
 */
func Info(c *gin.Context) {
	// 从上下文获取user
	user, _ := c.Get("user")
	response.Success(c, gin.H{"user": dto.ToUserDto(user.(model.User))}, "登陆成功") // 类型断言
	return
}

/**
 * @Description: 查询数据库中是否存在当前电话
 * @return bool
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
 * @Description: 随机产生指定长度的英文字符作为用户名
 * @return string
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
