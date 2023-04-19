// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luatable

import (
	"errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type marshalFunc func(encodingRun, protoreflect.Message) error

const GoogleProtobufParentPackage = "google.protobuf"
const GoogleProtobufTimestamp = "Timestamp"

// wellKnownTypeMarshaler returns a marshal function if the message type
// has specialized serialization behavior. It returns nil otherwise.
// it returns an error in case of types are not supported yet
func wellKnownTypesMarshaler(fullName protoreflect.FullName) (marshalFunc, error) {
	if fullName.Parent() == GoogleProtobufParentPackage {
		switch fullName.Name() {
		case GoogleProtobufTimestamp:
			// Timestamp can be converted to normal table because it can be converted to os.date and os.time with build in function in lua
			return nil, nil
		default:
			return nil, errors.New(string(fullName) + " is not supported yet.")
		}
	}
	return nil, nil
}
