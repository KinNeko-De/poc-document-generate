package function

import (
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KinNeko-De/restaurant-document-generate-function/internal/app/document-generate/document"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("HelloGoogleIAmANewCustomer", GenerateDocument)
}

type DocumentGenerator interface {
	GenerateDocument()
}

func GenerateDocument(ctx context.Context, e event.Event) error {
	return GenerateDocumentUsingGenerator(ctx, e, document.DocumentGenerator{})
}

func GenerateDocumentUsingGenerator(ctx context.Context, e event.Event, generator DocumentGenerator) error {
	// content := e.Data()
	// TODO parse protobuf
	generator.GenerateDocument()
	return nil
}
