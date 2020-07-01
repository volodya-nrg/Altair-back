package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

// GetTest - кастомная ф-ия тестирования. Просто для проверки тех или иных данных.
func GetTest(c *gin.Context) {
	c.JSON(200, nil)
}

// PrintMemUsage - вывод информации о памяти
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
