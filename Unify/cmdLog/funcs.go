package cmdLog

import (
	"log"
	"time"
)


func timeStamp() string {
	return time.Now().Local().Format(time.UnixDate) + " : "
}
func Printf(format string, args ...interface{}) {
	// format = timeStamp() + format
	log.Printf(format, args...)
}

func LogPrintDate(format string, args ...interface{}) {
	// log.Printf("%s ", TimeNow())
	log.Printf(format, args...)
	// fmt.Println()
}

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
