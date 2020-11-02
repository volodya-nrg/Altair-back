package server

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Db - главная перемення для DB
var Db *gorm.DB

// InitDB - ф-ия инициализации БД
func InitDB(user, password, host string, port uint, dbName string) error {
	var err error
	pattern := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"
	dbSettings := fmt.Sprintf(pattern, user, password, host, port, dbName)

	Db, err = gorm.Open(mysql.Open(dbSettings), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}
