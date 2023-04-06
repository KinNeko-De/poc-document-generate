package document

import (
	"log"
	"os/exec"
)

type DocumentGenerator struct {
}

func (_ DocumentGenerator) GenerateDocument() {
	document := "documents/invoice.tex"
	output := "/app/generated"
	outputParameter := "-output-directory=" + output
	cmd := exec.Command("lualatex", outputParameter, document)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("error executiong %v %v", cmd, err)
	}
	log.Println("Document generated.") // TODO make this debug
}
