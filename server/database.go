package server

import (
	"altair/pkg/logger"
	"fmt"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

func InitDB(user string, password string, host string, port uint, dbName string) error {
	pattern := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"
	dbSettings := fmt.Sprintf(pattern, user, password, host, port, dbName)
	var err error

	Db, err = gorm.Open("mysql", dbSettings)
	if err != nil {
		logger.Warning.Println(err)
		return err
	}

	return nil
}
