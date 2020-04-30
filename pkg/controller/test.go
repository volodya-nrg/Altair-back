package controller

import (
	"altair/api/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

func GetTest(c *gin.Context) {
	//PrintMemUsage()
	//timeStart := time.Now()
	//for i := 1000000; i >= 0; i-- {
	//	myTest()
	//}
	//PrintMemUsage()
	//fmt.Println("====>", time.Since(timeStart))
	//runtime.GC()
	//PrintMemUsage()

	//serviceCats := service.NewCatService()
	//cats, _ := serviceCats.GetCats()
	//catsTree := serviceCats.GetCatsAsTree(cats)
	//list := serviceCats.GetAncestors(catsTree, 479)

	//serviceAds := service.NewAdService()
	//ads, _ := serviceAds.GetAds("created_at desc")
	//for _, v := range ads {
	//	v.Description = helpers.Lorem(helpers.RandIntByRange(10, 300))
	//	v.Price = uint64(helpers.RandIntByRange(0, 999999999))
	//	v.Title = helpers.Lorem(helpers.RandIntByRange(10, 100))
	//	v.Youtube = "https://youtu.be/rX9-jFV8iz8"
	//	v.Ip = "127.0.0.1"
	//	serviceAds.Update(v, nil)
	//}

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
