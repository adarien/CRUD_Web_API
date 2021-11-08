package logs

import (
	"io"
	"log"
	"os"
)

var (
	// _INFO *logs.Logger
	ERROR *log.Logger
)

func initLog(errorHandle io.Writer) {

	// _INFO = logs.New(infoHandle, "INFO\t", logs.Ldate|logs.Ltime)
	ERROR = log.New(errorHandle, "ERROR\t", log.Ldate|log.Ltime)
}

func init() {
	logfile, err := os.Create("cmd/log.log")
	if err != nil {
		log.Fatalln(err)
	}
	initLog(logfile)
}
