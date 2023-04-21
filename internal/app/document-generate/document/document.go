package document

import (
	"github.com/google/uuid"
	protobuf "github.com/kinneko-de/test-api-contract/golang/kinnekode/protobuf"
	restaurantApi "github.com/kinneko-de/test-api-contract/golang/kinnekode/restaurant/document"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

type DocumentGenerator struct {
}

func (documentGenerator DocumentGenerator) GenerateDocument() {
	var request = CreateTestRequest()
	_ = ToLuaTable(request)

	currentDirectory := documentGenerator.getCurrentDirectory()
	luatexTemplateDirectory := path.Join(currentDirectory, "run")
	localDebugDirectory := path.Join(currentDirectory, "run", "generated")

	tmpDirectory := path.Join(currentDirectory, request.RequestId.Value)
	outputDirectory := path.Join(tmpDirectory, "generated")
	documentGenerator.createDirectoryForRun(outputDirectory)

	template := documentGenerator.GetTemplateName(request)
	templateFile := documentGenerator.CopyLuatexTemplate(luatexTemplateDirectory, template, tmpDirectory)

	documentGenerator.CreateDocumentInputData(luatexTemplateDirectory, template, tmpDirectory)

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

func (documentGenerator DocumentGenerator) GetTemplateName(request *restaurantApi.GenerateDocumentV1) string {
	var template string
	switch request.RequestedDocuments[0].Type.(type) {
	case *restaurantApi.GenerateDocumentV1_Document_Invoice:
		template = "invoice"
	default:
		log.Fatalf("Document %v not supported yet", request.RequestedDocuments[0].Type)
	}
	return template
}

func (documentGenerator DocumentGenerator) CopyLuatexTemplate(documentDirectory string, template string, tmpDirectory string) string {
	templateFile := template + ".tex"
	_, texErr := copyFile(path.Join(documentDirectory, templateFile), path.Join(tmpDirectory, templateFile))
	if texErr != nil {
		log.Fatalf("Can not copy tex file: %v", texErr)
	}
	return templateFile
}

func (documentGenerator DocumentGenerator) CreateDocumentInputData(luatexTemplateDirectory string, template string, tmpDirectory string) {
	inputDataFile := template + ".lua"
	_, luaErr := copyFile(path.Join(luatexTemplateDirectory, inputDataFile), path.Join(tmpDirectory, inputDataFile))
	if luaErr != nil {
		log.Fatalf("Can not copy lua file: %v", luaErr)
	}
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

func CreateTestRequest() *restaurantApi.GenerateDocumentV1 {
	randomRequestId := CreateRandomUuid()

	request := restaurantApi.GenerateDocumentV1{
		RequestId: randomRequestId,
		RequestedDocuments: []*restaurantApi.GenerateDocumentV1_Document{
			{
				Type: &restaurantApi.GenerateDocumentV1_Document_Invoice{
					Invoice: &restaurantApi.GenerateDocumentV1_Document_InvoiceV1{
						DeliveredOn:  timestamppb.New(time.Date(2020, time.April, 13, 0, 0, 0, 0, time.UTC)),
						CurrencyCode: "EUR",
						Recipient: &restaurantApi.GenerateDocumentV1_Document_InvoiceV1_Recipient{
							Name:     "Max Mustermann",
							Street:   "Musterstraße 17",
							City:     "Musterstadt",
							PostCode: "12345",
							Country:  "DE",
						},
						Items: []*restaurantApi.GenerateDocumentV1_Document_InvoiceV1_Item{
							{
								Description: "vfdsdsfdsfdsfs fdsfdskfdsk fdskfk fkwef kefkwekfe\\r\\nANS 23054303053",
								Quantity:    2,
								NetAmount:   &protobuf.Decimal{Value: "3.35"},
								Taxation:    &protobuf.Decimal{Value: "19"},
								TotalAmount: &protobuf.Decimal{Value: "3.99"},
								Sum:         &protobuf.Decimal{Value: "7.98"},
							},
							{
								Description: "vf ds dsf dsf dsfs fds fd skf dsk\\r\\nANS 606406540",
								Quantity:    1,
								NetAmount:   &protobuf.Decimal{Value: "9.07"},
								Taxation:    &protobuf.Decimal{Value: "19"},
								TotalAmount: &protobuf.Decimal{Value: "10.79"},
								Sum:         &protobuf.Decimal{Value: "10.79"},
							},
							{
								Description: "Versandkosten",
								Quantity:    1,
								NetAmount:   &protobuf.Decimal{Value: "0.00"},
								Taxation:    &protobuf.Decimal{Value: "0"},
								TotalAmount: &protobuf.Decimal{Value: "0.00"},
								Sum:         &protobuf.Decimal{Value: "0.00"},
							},
						},
					},
				},
				OutputFormats: []restaurantApi.GenerateDocumentV1_Document_OutputFormat{
					restaurantApi.GenerateDocumentV1_Document_OUTPUT_FORMAT_PDF,
				},
			},
			{},
		},
	}

	return &request
}

func CreateRandomUuid() *protobuf.Uuid {
	id, uuidErr := uuid.NewUUID()
	if uuidErr != nil {
		log.Fatalf("error generating google uuid: %v", uuidErr)
	}
	randomRequestId, protobufErr := protobuf.ToProtobuf(id)
	if protobufErr != nil {
		log.Fatalf("error generating protobuf uuid: %v", protobufErr)
	}
	return randomRequestId
}

func ToLuaTable(_ *restaurantApi.GenerateDocumentV1) []byte {
	return []byte{}
}
