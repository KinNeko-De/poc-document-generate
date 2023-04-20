package luatable

import "google.golang.org/protobuf/reflect/protoreflect"

type converter func(encodingRun, protoreflect.Message) error

// UserConverter defines how user specific protobuf message are converted to lua tables.
// UserConverter returns a converter function if the message type should be handled by the converter.
// UserConverter returns nil if it does not want to handle the message
type UserConverter interface {
	Handle(fullName protoreflect.FullName) (converter, error)
}
