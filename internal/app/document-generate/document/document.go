package document

import (
	"github.com/KinNeko-De/restaurant-document-generate-function/encoding/luatable"
	restaurantApi "github.com/kinneko-de/test-api-contract/golang/kinnekode/restaurant/document"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type DocumentGenerator struct {
}

func (documentGenerator DocumentGenerator) GenerateDocument(request *restaurantApi.GenerateDocumentV1) {
	currentDirectory := documentGenerator.getCurrentDirectory()
	luatexTemplateDirectory := path.Join(currentDirectory, "run")
	localDebugDirectory := path.Join(currentDirectory, "run", "generated")

	tmpDirectory := path.Join(currentDirectory, request.RequestId.Value)
	outputDirectory := path.Join(tmpDirectory, "generated")
	documentGenerator.createDirectoryForRun(outputDirectory)

	template, message := documentGenerator.GetTemplateName(request)
	templateFile := documentGenerator.CopyLuatexTemplate(luatexTemplateDirectory, template, tmpDirectory)

	documentInputData := ToLuaTable(message)
	documentGenerator.CreateDocumentInputData(template, tmpDirectory, documentInputData)

	documentGenerator.ExecuteLuaLatex(outputDirectory, templateFile, tmpDirectory, localDebugDirectory)
	log.Println("Document generated.") // TODO make this debug
}

func (documentGenerator DocumentGenerator) ExecuteLuaLatex(outputDirectory string, templateFile string, tmpDirectory string, localDebugDirectory string) {
	cmd, commandError := documentGenerator.runCommand(outputDirectory, templateFile, tmpDirectory)

	if commandError != nil {
		documentGenerator.SecureLogFileForLocalDebug(outputDirectory, templateFile, localDebugDirectory)
		log.Fatalf("error executing %v %v", cmd, commandError)
	} else {
		documentGenerator.SecurePdfForLocalDebug(outputDirectory, localDebugDirectory)
	}
}

func (documentGenerator DocumentGenerator) SecurePdfForLocalDebug(outputDirectory string, localDebugDirectory string) {
	_, pdfErr := copyFile(path.Join(outputDirectory, "invoice.pdf"), path.Join(localDebugDirectory, "invoice.pdf"))
	if pdfErr != nil {
		log.Fatalf("error coping pdf file %v", pdfErr)
	}
}

func (documentGenerator DocumentGenerator) SecureLogFileForLocalDebug(outputDirectory string, templateFile string, localDebugDirectory string) {
	luaLatexLog := templateFile[:len(templateFile)-len(filepath.Ext(templateFile))] + ".log"
	source := path.Join(outputDirectory, luaLatexLog)
	destination := path.Join(localDebugDirectory, luaLatexLog)
	_, logCopyErr := copyFile(source, destination)
	if logCopyErr != nil {
		log.Fatalf("error coping latex log from %v to %v: %v", source, destination, logCopyErr)
	}
}

func (documentGenerator DocumentGenerator) runCommand(outputDirectory string, templateFile string, tmpDirectory string) (*exec.Cmd, error) {
	outputParameter := "-output-directory=" + outputDirectory
	cmd := exec.Command("lualatex", outputParameter, templateFile)
	cmd.Dir = tmpDirectory
	commandError := cmd.Run()
	return cmd, commandError
}

func (documentGenerator DocumentGenerator) GetTemplateName(request *restaurantApi.GenerateDocumentV1) (string, proto.Message) {
	var template string
	var message proto.Message
	switch request.RequestedDocuments[0].Type.(type) {
	case *restaurantApi.GenerateDocumentV1_Document_Invoice:
		template = "invoice"
		message = request.RequestedDocuments[0].GetInvoice()
	default:
		log.Fatalf("Document %v not supported yet", request.RequestedDocuments[0].Type)
	}
	return template, message
}

func (documentGenerator DocumentGenerator) CopyLuatexTemplate(documentDirectory string, template string, tmpDirectory string) string {
	templateFile := template + ".tex"
	_, texErr := copyFile(path.Join(documentDirectory, templateFile), path.Join(tmpDirectory, templateFile))
	if texErr != nil {
		log.Fatalf("Can not copy tex file: %v", texErr)
	}
	return templateFile
}

func (documentGenerator DocumentGenerator) CreateDocumentInputData(template string, tmpDirectory string, inputData []byte) {
	inputDataFile := template + ".lua"
	file, err := os.Create(path.Join(tmpDirectory, inputDataFile))
	if err != nil {
		log.Fatalf("Error creating input data: %v", err)
	}
	file.WriteString("local ")
	file.Write(inputData)
	// TODO Make name of InvoiceV1 flexible
	tableAssign := "return {" + template + " = InvoiceV1 }"
	file.WriteString(tableAssign)
	file.Close()
}

func (_ DocumentGenerator) createDirectoryForRun(outputDirectory string) {
	mkDirError := os.MkdirAll(outputDirectory, os.ModeExclusive)
	if mkDirError != nil {
		log.Fatalf("Can not create output directory: %v", mkDirError)
	}
}

func (_ DocumentGenerator) getCurrentDirectory() string {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("error get current directory: %v", err)
	}
	return currentDirectory
}

func copyFile(src, dst string) (int64, error) {
	source, openError := os.Open(src)
	if openError != nil {
		return 0, openError
	}
	defer source.Close()

	destination, createError := os.Create(dst)
	if createError != nil {
		return 0, createError
	}
	defer destination.Close()
	nBytes, copyError := io.Copy(destination, source)
	return nBytes, copyError
}

func ToLuaTable(m proto.Message) []byte {
	opt := luatable.MarshalOptions{Multiline: true, UserConverters: []luatable.UserConverter{luatable.KinnekodeProtobuf{}}}
	luaTable, err := opt.Marshal(m)
	if err != nil {
		log.Fatalf("Error converting protobuf message to luat table: %v", err)
	}
	return luaTable
}
