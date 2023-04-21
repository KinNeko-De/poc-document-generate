package function

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	restaurantApi "github.com/kinneko-de/test-api-contract/golang/kinnekode/restaurant/document"
	"log"
	"testing"
)

type MockGenerator struct {
}

func (_ MockGenerator) GenerateDocument(*restaurantApi.GenerateDocumentV1) {
	log.Println("Dummy implementation that does nothing.")
}

func TestGenerateDocument(t *testing.T) {
	data := []byte{} // TODO getProtobufMessageHere
	event := event.Event{DataEncoded: data}
	GenerateDocumentUsingGenerator(context.Background(), event, MockGenerator{})
	// TODO Asserts
}
