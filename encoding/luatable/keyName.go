package luatable

import "google.golang.org/protobuf/reflect/protoreflect"

type keyName interface {
	keyName(protoreflect.FieldDescriptor) string
}

type jsonName struct {
}

func (j jsonName) keyName(fieldDescriptor protoreflect.FieldDescriptor) string {
	return fieldDescriptor.JSONName()
}

type protobufName struct {
}

func (p protobufName) keyName(fieldDescriptor protoreflect.FieldDescriptor) string {
	return fieldDescriptor.TextName()
}
