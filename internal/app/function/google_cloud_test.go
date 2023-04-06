package function

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"log"
	"testing"
)

type MockGenerator struct {
}

func (_ MockGenerator) GenerateDocument() {
	log.Println("Dummy implementation that does nothing.")
}

func TestGenerateDocument(t *testing.T) {
	data := []byte{} // TODO getProtobufMessageHere
	event := event.Event{DataEncoded: data}
	GenerateDocumentUsingGenerator(context.Background(), event, MockGenerator{})
	// TODO Asserts
}
