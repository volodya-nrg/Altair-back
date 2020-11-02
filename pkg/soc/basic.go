package soc

import (
	"altair/pkg/manager"
	"fmt"
)

// Socer - общий интерфейс для соц. сетей (аунтификация)
type Socer interface {
	GetAccessToken(string) error
	GetUserInfo() (interface{}, error)
	checkAccessToken() error
}

// CommonHandler - общий обработчик соц. классов
func CommonHandler(opt string, obj Socer) (CommonUserInfo, error) {
	commonUserInfo := CommonUserInfo{}

	if err := obj.GetAccessToken("https://www.altair.uz/login"); err != nil {
		return commonUserInfo, err
	}

	tmpUserInfo, err := obj.GetUserInfo()
	if err != nil {
		return commonUserInfo, err
	}

	switch opt {
	case "vk":
		data, ok := tmpUserInfo.(VkUserInfo)
		if !ok {
			return commonUserInfo, manager.ErrInConvertType
		}

		commonUserInfo.Email = fmt.Sprintf("id%d@vk.com", data.ID)
		commonUserInfo.Name = data.FirstName
	case "ok":
		data, ok := tmpUserInfo.(OkUserInfo)
		if !ok {
			return commonUserInfo, manager.ErrInConvertType
		}

		commonUserInfo.Email = fmt.Sprintf("id%s@ok.ru", data.UID)
		commonUserInfo.Name = data.FirstName
	case "fb":
		data, ok := tmpUserInfo.(FbUserInfo)
		if !ok {
			return commonUserInfo, manager.ErrInConvertType
		}

		commonUserInfo.Email = fmt.Sprintf("id%s@facebook.com", data.ID)
		commonUserInfo.Name = data.Name
	case "ggl":
		data, ok := tmpUserInfo.(GglUserInfo)
		if !ok {
			return commonUserInfo, manager.ErrInConvertType
		}

		commonUserInfo.Email = fmt.Sprintf("id%s@google.com", data.ID)
		commonUserInfo.Name = data.Name
	default:
		return commonUserInfo, manager.ErrUndefinedOptSoc
	}

	return commonUserInfo, nil
}

// CommonUserInfo - обобщающая структура данных пользователя
type CommonUserInfo struct {
	Email string
	Name  string
}
