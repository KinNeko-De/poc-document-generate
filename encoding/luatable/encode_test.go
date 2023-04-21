package luatable_test

import (
	"github.com/KinNeko-De/restaurant-document-generate-function/encoding/luatable"
	"github.com/google/go-cmp/cmp"
	"github.com/kinneko-de/test-api-contract/golang/kinnekode/protobuf"
	restaurantApi "github.com/kinneko-de/test-api-contract/golang/kinnekode/restaurant/document"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestFormatPositive(t *testing.T) {
	tests := []struct {
		desc     string
		input    proto.Message
		option   luatable.MarshalOptions
		expected string
	}{
		{
			desc: "Message with string",
			input: &restaurantApi.GenerateDocumentV1_Document_InvoiceV1_Recipient{
				Name:     "Max Mustermann",
				Street:   "Musterstraße 17",
				City:     "Musterstadt",
				PostCode: "12345",
				Country:  "DE",
			},
			expected: "Recipient={name=\"Max Mustermann\",street=\"Musterstraße 17\",city=\"Musterstadt\",postCode=\"12345\",country=\"DE\"}",
		},
		{
			desc: "long test",
			input: &restaurantApi.GenerateDocumentV1_Document_InvoiceV1{
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
						Description: "vfd % \\r\\nANS 23054303053",
						Quantity:    2,
						NetAmount:   &protobuf.Decimal{Value: "3.35"},
						Taxation:    &protobuf.Decimal{Value: "19"},
						TotalAmount: &protobuf.Decimal{Value: "3.99"},
						Sum:         &protobuf.Decimal{Value: "7.98"},
					},
					{
						Description: "Abv djefk\\r\\nANS 606406540",
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

			option:   luatable.MarshalOptions{Multiline: true, UserConverters: []luatable.UserConverter{luatable.KinnekodeProtobuf{}}},
			expected: "InvoiceV1 = {\n  deliveredOn = {\n    seconds = 1586736000\n  },\n  currencyCode = \"EUR\",\n  recipient = {\n    name = \"Max Mustermann\",\n    street = \"Musterstraße 17\",\n    city = \"Musterstadt\",\n    postCode = \"12345\",\n    country = \"DE\"\n  },\n  items = {\n    [1] = {\n      description = \"vfd \\\\% \\\\\\\\ANS 23054303053\",\n      quantity = 2,\n      netAmount = 3.35,\n      taxation = 19,\n      totalAmount = 3.99,\n      sum = 7.98\n    },\n    [2] = {\n      description = \"Abv djefk\\\\\\\\ANS 606406540\",\n      quantity = 1,\n      netAmount = 9.07,\n      taxation = 19,\n      totalAmount = 10.79,\n      sum = 10.79\n    },\n    [3] = {\n      description = \"Versandkosten\",\n      quantity = 1,\n      netAmount = 0.00,\n      taxation = 0,\n      totalAmount = 0.00,\n      sum = 0.00\n    }\n  }\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			b, err := tt.option.Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal() returned error: %v\n", err)
			}
			actual := string(b)
			if actual != tt.expected {
				t.Errorf("Marshal()\n<actual>\n%v\n<expected>\n%v\n", actual, tt.expected)
				if diff := cmp.Diff(tt.expected, actual); diff != "" {
					t.Errorf("Marshal() diff -expected +actual\n%v\n", diff)
				}
			}
		})
	}
}
