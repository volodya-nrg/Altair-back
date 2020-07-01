package server

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// Db - главная перемення для DB
var Db *gorm.DB

// InitDB - ф-ия инициализации БД
func InitDB(user string, password string, host string, port uint, dbName string) error {
	var err error
	pattern := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"
	dbSettings := fmt.Sprintf(pattern, user, password, host, port, dbName)

	Db, err = gorm.Open("mysql", dbSettings)
	if err != nil {
		return err
	}

	return nil
}
