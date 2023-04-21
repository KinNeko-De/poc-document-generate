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
	"time"
)

type DocumentGenerator struct {
}

func (_ DocumentGenerator) GenerateDocument() {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	documentFolder := path.Join(currentDirectory, "run")
	localDebugFolder := path.Join(currentDirectory, "run", "generated")

	var request = CreateTestRequest()
	_ = ToLuaTable(request)

	tmpDirectory := path.Join(currentDirectory, request.RequestId.Value)
	output := path.Join(tmpDirectory, "generated")
	mkDirError := os.MkdirAll(output, os.ModeExclusive)
	if mkDirError != nil {
		log.Fatalf("Can not create output directory: %v", mkDirError)
	}
	document := "invoice.tex"
	_, texErr := copyFile(path.Join(documentFolder, document), path.Join(tmpDirectory, document))
	if texErr != nil {
		log.Fatalf("Can not copy tex file: %v", mkDirError)
	}
	_, luaErr := copyFile(path.Join(documentFolder, "invoice.lua"), path.Join(tmpDirectory, "invoice.lua"))
	if luaErr != nil {
		log.Fatalf("Can not copy lua file: %v", mkDirError)
	}
	outputParameter := "-output-directory=" + output
	cmd := exec.Command("lualatex", outputParameter, document)
	cmd.Dir = tmpDirectory
	commandError := cmd.Run()

	if commandError != nil {
		source := path.Join(output, "invoice.log")
		destination := path.Join(localDebugFolder, "invoice.log")
		_, logCopyErr := copyFile(source, destination)
		if logCopyErr != nil {
			log.Fatalf("error coping latex log from %v to %v: %v", source, destination, logCopyErr)
		}
		log.Println("going to sleep for debug")
		time.Sleep(time.Minute * 30)
		log.Fatalf("error executing %v %v", cmd, commandError)
	} else {
		_, pdfErr := copyFile(path.Join(output, "invoice.pdf"), path.Join(localDebugFolder, "invoice.pdf"))
		if pdfErr != nil {
			log.Fatalf("error coping pdf file %v", pdfErr)
		}
	}
	log.Println("Document generated.") // TODO make this debug
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
