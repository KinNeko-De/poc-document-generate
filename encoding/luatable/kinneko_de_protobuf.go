// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luatable

import (
	"errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const KinNekoDeProtobufParentPackage = "kinnekode.protobuf"
const KinNekoDeProtobufDecimal = "Decimal"

type KinnekodeProtobuf struct {
}

func (KinnekodeProtobuf) Handle(fullName protoreflect.FullName) (converter, error) {
	if fullName.Parent() == KinNekoDeProtobufParentPackage {
		switch fullName.Name() {
		case KinNekoDeProtobufDecimal:
			return encodingRun.convertDecimal, nil
		default:
			return nil, errors.New(string(fullName) + " is not supported yet")
		}
	}
	return nil, nil
}

func (e encodingRun) convertDecimal(m protoreflect.Message) error {
	fd := m.Descriptor().Fields().ByNumber(1)
	val := m.Get(fd)
	e.WriteNumber(val.String())
	return nil
}
