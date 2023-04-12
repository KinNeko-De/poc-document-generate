package operation

import (
	"log"
	"os"
)

func UseLogFileInGenerated() *os.File {
	f, err := os.OpenFile("/app/run/generated/testlog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	log.SetOutput(f)
	return f
}

func CloseLogFile(logfile *os.File) {
	err := logfile.Close()
	if err != nil {
		log.Fatalf("error closing log : %v", err)
	}
}
