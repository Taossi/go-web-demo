package middleware

import (
	"gin-gorm/common"
	"gin-gorm/model"
	"gin-gorm/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/**
 * @Description: 检验token的中间件，保护路由
				 用户获得token后需要携带token访问其他接口 需要进行权限校验
 * @return gin.HandlerFunc
*/
func AuthMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头, 如 Authorization: Bearer mF_9.B5f-4.1JqM
		// 获得authorization header
		authString := context.GetHeader("Authorization")
		// 验证token格式
		if authString == "" || !strings.HasPrefix(authString, "Bearer ") {
			response.Response(context, http.StatusUnauthorized, 401, nil, "token权限错误")
			context.Abort()
			return
		}
		authString = authString[7:]

		// 解析 校验token
		claims, err := common.ParseToken(authString)
		if err != nil {
			response.Response(context, http.StatusUnauthorized, 401, nil, "token权限错误")
			context.Abort()
			return
		}

		// 验证通过后获取userID
		userId := claims.UserId
		DB := common.InitDB()
		user := model.User{}
		DB.First(&user, userId)

		// user不存在
		if userId == 0 {
			response.Response(context, http.StatusUnauthorized, 401, nil, "token权限错误")
			context.Abort()
			return
		}

		// 将当前请求的user信息保存到请求的上下文context上
		context.Set("user", user)
		context.Next()
	}
}
