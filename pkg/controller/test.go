package controller

import (
	"altair/api/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

func GetTest(c *gin.Context) {
	//logger.Info.Println("enter in GetTest")
	//defer logger.Info.Println("exit in mGetTestyfn")
	//
	//go func() {
	//	logger.Info.Println("enter in gor")
	//	defer logger.Info.Println("exit in gor")
	//	time.Sleep(time.Second * 3)
	//	logger.Info.Println("do gor")
	//}()

	//PrintMemUsage()
	//timeStart := time.Now()
	//for i := 1000000; i >= 0; i-- {
	//	myTest()
	//}
	//PrintMemUsage()
	//fmt.Println("====>", time.Since(timeStart))
	//runtime.GC()
	//PrintMemUsage()

	c.JSON(200, nil)
}
func myTest() {
	res := response.Result{}
	res.Status = 200
	res.Err = errors.New("err test")
	res.Data = "test"
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
