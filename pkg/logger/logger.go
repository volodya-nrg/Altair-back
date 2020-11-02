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

// Init - ф-ия инициализации логгера
func Init(infoHandle, warningHandle, errorHandle io.Writer, isDebugMode bool) {
	prefixInfo, prefixWarning, prefixError := "INFO: ", "WARNING: ", "ERROR: "
	iInfo := log.LstdFlags | log.Llongfile
	iWarning := log.LstdFlags | log.Llongfile
	iError := log.LstdFlags | log.Llongfile

	if isDebugMode {
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		prefixInfo = green(prefixInfo)
		prefixWarning = yellow(prefixWarning)
		prefixError = red(prefixError)
		iInfo, iWarning, iError = log.Lshortfile, log.Lshortfile, log.Lshortfile
	}

	Info = log.New(infoHandle, prefixInfo, iInfo)
	Warning = log.New(warningHandle, prefixWarning, iWarning)
	Error = log.New(errorHandle, prefixError, iError)
}
