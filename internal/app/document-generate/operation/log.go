package operation

import (
	"log"
	"os"
)

func UseLogFileInGenerated() *os.File {
	f, err := os.OpenFile("/app/generated/testlog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}

func CloseLogFile(logfile *os.File) {
	logfile.Close()
}
