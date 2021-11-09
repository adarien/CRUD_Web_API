package logs

import (
	"io"
	"log"
	"os"
)

var (
	INFO  *log.Logger
	ERROR *log.Logger
)

func initLog(infoHandle io.Writer, errorHandle io.Writer) {

	INFO = log.New(infoHandle, "INFO\t", log.Ldate|log.Ltime)
	ERROR = log.New(errorHandle, "ERROR\t", log.Ldate|log.Ltime)
}

func init() {
	logfile, err := os.Create("logs/log.log")
	if err != nil {
		log.Fatalln(err)
	}
	initLog(logfile, logfile)
}
