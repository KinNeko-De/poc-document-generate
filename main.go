package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	f, ferr := os.OpenFile("/app/generated/testlog2.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if ferr != nil {
		log.Fatalf("error opening file: %v", ferr)
	}
	defer f.Close()
	log.SetOutput(f)

	document := "documents/invoice.tex"
	output := "/app/generated"
	outputParameter := "-output-directory=" + output
	cmd := exec.Command("lualatex", outputParameter, document)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("error executiong %v %v", cmd, err)
	}
	log.Println("Document generated.")
}
