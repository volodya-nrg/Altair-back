package logger

import (
	"github.com/fatih/color"
	"io"
	"log"
)

var (
	// Info - информирующий логгер
	Info *log.Logger
	// Warning - предупреждающий логгер
	Warning *log.Logger
	// Error - логгер ошибки
	Error *log.Logger
)

// Init - ф-ия уставновки логгера
func Init(infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	Info = log.New(infoHandle, green("INFO: "), log.Lshortfile)
	Warning = log.New(warningHandle, yellow("WARNING: "), log.Lshortfile)
	Error = log.New(errorHandle, red("ERROR: "), log.Lshortfile)
}

//f, err := os.OpenFile("./errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//if err != nil {
//	log.Fatalf("error opening file: %v", err)
//}
//defer f.Close()
//
//log.SetOutput(f)
//log.Println("This is a test log entry")
