package controller

import (
	"github.com/gin-gonic/gin"
)

// GetTest - кастомная ф-ия тестирования. Просто для проверки тех или иных данных.
func GetTest(c *gin.Context) {
	//var data = []string{"1", "2", "3", "4"}
	//ch := make(chan struct{}, 0)
	//
	//for _, v := range data {
	//	fmt.Println("==iterator==")
	//	go func(vSrc string) {
	//		myInit(vSrc)
	//		<-ch
	//	}(v)
	//	ch <- struct{}{}
	//}

	// data := map[string]string{} // изначально передается по ссылке.
	// myMap(data)
	// fmt.Printf("--->%#v - %s\n", data, reflect.ValueOf(data).Kind())

	// data := A{B: "testStruct"} // передается по значению
	// changeStruct(data)
	// fmt.Printf("--->%#v - %s\n", data, reflect.ValueOf(data).Kind())

	//data := []string{"1", "2", "3", "4"} // передается по значению
	//changeSlice(data)
	//fmt.Printf("--->%#v - %s\n", data, reflect.ValueOf(data).Kind())

	//data := [5]int{1000, 2, 3, 17, 50} // передается по значению
	//changeArray(data)
	//fmt.Printf("--->%#v - %s\n", data, reflect.ValueOf(data).Kind())

	//fmt.Println("==end==")
	//fmt.Println(a("([{}]){}"))
	c.JSON(200, nil)
}

//func a(str string) bool {
//	var stack = make([]string, 0)
//	var count1 uint64
//	var count2 uint64
//	var count3 uint64
//
//	for _, v := range str {
//		symbol := string(v)
//
//		if symbol == "(" {
//			stack = append(stack, symbol)
//			count1++
//		} else if symbol == "[" {
//			stack = append(stack, symbol)
//			count2++
//		} else if symbol == "{" {
//			stack = append(stack, symbol)
//			count3++
//		}
//
//		if symbol == ")" {
//			count1--
//
//			if len(stack) > 0 {
//				s := stack[len(stack)-1]
//				stack = stack[:len(stack)-1]
//
//				if s != "(" {
//					return false
//				}
//			}
//
//		} else if symbol == "]" {
//			count2--
//
//			if len(stack) > 0 {
//				s := stack[len(stack)-1]
//				stack = stack[:len(stack)-1]
//
//				if s != "[" {
//					return false
//				}
//			}
//
//		} else if symbol == "}" {
//			count3--
//
//			if len(stack) > 0 {
//				s := stack[len(stack)-1]
//				stack = stack[:len(stack)-1]
//
//				if s != "{" {
//					return false
//				}
//			}
//		}
//	}
//
//	fmt.Printf("%v", stack)
//
//	return len(stack) == 0 && count1 == 0 && count2 == 0 && count3 == 0
//}

// PrintMemUsage - вывод информации о памяти
//func PrintMemUsage() {
//	var m runtime.MemStats
//	runtime.ReadMemStats(&m)
//	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
//	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
//	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
//	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
//	fmt.Printf("\tNumGC = %v\n", m.NumGC)
//}
//func bToMb(b uint64) uint64 {
//	return b / 1024 / 1024
//}
//func myInit(s string) {
//	time.Sleep(time.Second * 2)
//	fmt.Println("==" + s + "==")
//}
//func myMap(x map[string]string) {
//	x["asd"] = "test"
//}
//func changeStruct(x A) {
//	x.B = "changeStruct"
//}
//func changeSlice(x []string) {
//	x = append(x, "5")
//}
//func changeArray(arr [5]int) {
//	arr[len(arr)-1] = 10
//}
//type A struct {
//	B string
//}
