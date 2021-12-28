package common

import (
	"errors"
	"gin-gorm/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// 定义jwt过期时间
const TokenExpireDuration = time.Hour * 7 * 24

// 定义key secret
var jwtKey = []byte("secret")

/**
 * @Description: 发放生成token
 * @return string, error
 */
func ReleaseToken(user model.User) (string, error) {
	claims := Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),                          // 发放时间
			Issuer:    "fht",                                      // 签发人
			Subject:   "user token",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

/**
 * @Description: 根据传入的token值获取到Claims对象信息
 * @param tokenString
 */
func ParseToken(tokenString string) (*Claims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
