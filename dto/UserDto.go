package dto

import "gin-gorm/model"

// 定义数据传输对象 只返回name, phone
type UserDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:      user.Name,
		Telephone: user.Phone,
	}
}
