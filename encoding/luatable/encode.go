package luatable

import (
	"errors"
	"fmt"
	"github.com/KinNeko-De/restaurant-document-generate-function/internal/encoding/luatable"
	"github.com/KinNeko-De/restaurant-document-generate-function/internal/order"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const defaultIndent = "  "

// Format formats the message as a lua table using default options.
// Supports only proto3 messages
func Format(m proto.Message) string {
	return MarshalOptions{}.Format(m)
}

// Marshal convert the given proto.Message into a lua table using default options.
// Supports only proto3 messages
func Marshal(m proto.Message) ([]byte, error) {
	return MarshalOptions{}.Marshal(m)
}

type (
	// MarshalOptions configures how the lua table is created
	MarshalOptions struct {
		// Multiline specifies whether the marshaler should format the output in
		// indented-form with every textual element on a new line.
		// If Indent is an empty string, then the default indent is chosen.
		Multiline bool

		// Indent specifies the set of indentation characters to use in a multiline
		// formatted output such that every entry is preceded by Indent and
		// terminated by a newline. If non-empty, then Multiline is treated as true.
		// Indent can only be composed of space or tab characters.
		Indent string

		// How name of the key is defined
		// if not set, the default are jsonName
		// remarks: does not work for the initial message because we do not have a fieldDescriptor there
		KeyName keyName

		UserConverters []UserConverter

		// UseEnumNumbers emits enum values as numbers.
		UseEnumNumbers bool

		// Resolver is used for looking up types when expanding google.protobuf.Any
		// messages. If nil, this defaults to using protoregistry.GlobalTypes.
		Resolver interface {
			protoregistry.ExtensionTypeResolver
			protoregistry.MessageTypeResolver
		}
	}
)

// Format formats the message as a lua table using default options.
func (o MarshalOptions) Format(m proto.Message) string {
	if m == nil {
		return "{}" // message is nil
	}
	if !m.ProtoReflect().IsValid() {
		return "The message has a invalid format."
	}
	b, _ := o.Marshal(m)
	return string(b)
}

// Marshal convert the given proto.Message into a lua table using the options.
func (o MarshalOptions) Marshal(m proto.Message) ([]byte, error) {
	return o.marshal(m)
}

// marshal is a centralized function that all marshal operations go through with the initial message.
func (o MarshalOptions) marshal(m proto.Message) ([]byte, error) {
	if m == nil {
		return nil, errors.New("message can not be nil")
	}

	SetDefaults(&o)

	encoder, err := luatable.NewEncoder(o.Indent)
	if err != nil {
		return nil, err
	}

	enc := encodingRun{encoder, o}

	bytes, err2 := marshalRootMessage(m.ProtoReflect(), enc)
	if err2 != nil {
		return bytes, err2
	}

	return enc.Bytes(), nil
}

func SetDefaults(o *MarshalOptions) {
	if o.Multiline && o.Indent == "" {
		o.Indent = defaultIndent
	}
	if o.Resolver == nil {
		o.Resolver = protoregistry.GlobalTypes
	}
	if o.KeyName == nil {
		o.KeyName = jsonName{}
	}
}

func marshalRootMessage(m protoreflect.Message, enc encodingRun) ([]byte, error) {
	// The json name is not populated, so the Protobuf name is used here
	enc.WriteKey(string(m.Descriptor().Name()))

	if err := enc.marshalMessage(m); err != nil {
		return nil, err
	}
	return nil, nil
}

type encodingRun struct {
	*luatable.Encoder
	opts MarshalOptions
}

// marshalMessage marshals the message and fields in the given protoreflect.Message.
// inject the table name as property of the field descriptor if there is one, otherwise set on. Json names are not populated in the Descriptor :/
func (e encodingRun) marshalMessage(m protoreflect.Message) error {
	for _, userConverter := range e.opts.UserConverters {
		myUserConverter, unsupportedTypeError := userConverter.Handle(m.Descriptor().FullName())
		if unsupportedTypeError != nil {
			return unsupportedTypeError
		}
		if myUserConverter != nil {
			return myUserConverter(e, m)
		}
	}

	specialMarshal, unsupportedTypeError := wellKnownTypesMarshaler(m.Descriptor().FullName())
	if unsupportedTypeError != nil {
		return unsupportedTypeError
	}
	if specialMarshal != nil {
		return specialMarshal(e, m)
	}

	e.StartObject()
	defer e.EndObject()

	var fields order.FieldRanger = m

	var err error
	order.RangeFields(fields, order.IndexNameFieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		name := e.opts.KeyName.keyName(fd)

		if err = e.WriteKey(name); err != nil {
			return false
		}
		if err = e.marshalValue(v, fd); err != nil {
			return false
		}
		return true
	})
	return err
}

// marshalValue marshals the given protoreflect.Value.
func (e encodingRun) marshalValue(val protoreflect.Value, fd protoreflect.FieldDescriptor) error {
	switch {
	case fd.IsList():
		return e.marshalList(val.List(), fd)
	case fd.IsMap():
		return e.marshalMap(val.Map(), fd)
	default:
		return e.marshalSingular(val, fd)
	}
}

// marshalSingular marshals the given non-repeated field value. This includes
// all scalar types, enums, messages, and groups.
func (e encodingRun) marshalSingular(val protoreflect.Value, fd protoreflect.FieldDescriptor) error {
	if !val.IsValid() {
		return nil
	}

	switch kind := fd.Kind(); kind {
	case protoreflect.BoolKind:
		e.WriteBool(val.Bool())

	case protoreflect.StringKind:
		if err := e.WriteString(val.String()); err != nil {
			return err
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		e.WriteInt(val.Int())
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		e.WriteUint(val.Uint())
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		// 64-bit integers are written out as JSON string.
		e.WriteNumber(val.String())

	case protoreflect.FloatKind:
		errors.New("Floats are not supported yet.")

	case protoreflect.DoubleKind:
		errors.New("Floats are not supported yet.")

	case protoreflect.BytesKind:
		errors.New("Bytes are not supported yet.")

	case protoreflect.EnumKind:
		errors.New("Enum are not supported yet.")

	case protoreflect.MessageKind, protoreflect.GroupKind:
		if err := e.marshalMessage(val.Message()); err != nil {
			return err
		}

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	}
	return nil
}

// marshalList marshals the given protoreflect.List.
func (e encodingRun) marshalList(list protoreflect.List, fd protoreflect.FieldDescriptor) error {

	e.StartArray()
	defer e.EndArray()

	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		e.Encoder.WriteIndexedList(i + 1)
		if err := e.marshalSingular(item, fd); err != nil {
			return err
		}
	}
	return nil
}

// marshalMap marshals given protoreflect.Map.
func (e encodingRun) marshalMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor) error {
	e.StartObject()
	defer e.EndObject()

	var err error
	order.RangeEntries(mmap, order.GenericKeyOrder, func(k protoreflect.MapKey, v protoreflect.Value) bool {
		if err = e.WriteKey(k.String()); err != nil {
			return false
		}
		if err = e.marshalSingular(v, fd.MapValue()); err != nil {
			return false
		}
		return true
	})
	return err
}
