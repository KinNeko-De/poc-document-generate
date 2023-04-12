package document

import (
	"log"
	"os/exec"
)

type DocumentGenerator struct {
}

func (_ DocumentGenerator) GenerateDocument() {
	document := "invoice.tex"
	output := "./generated"
	outputParameter := "-output-directory=" + output
	cmd := exec.Command("lualatex", outputParameter, document)
	cmd.Dir = "/app/run"
	err := cmd.Run()
	if err != nil {
		log.Fatalf("error executiong %v %v", cmd, err)
	}
	log.Println("Document generated.") // TODO make this debug
}
