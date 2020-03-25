package controller

import (
	"altair/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func GetTest(c *gin.Context) {
	logger.Info.Println("enter in GetTest")
	defer logger.Info.Println("exit in mGetTestyfn")

	go func() {
		logger.Info.Println("enter in gor")
		defer logger.Info.Println("exit in gor")
		time.Sleep(time.Second * 3)
		logger.Info.Println("do gor")
	}()

	c.JSON(200, nil)
}
