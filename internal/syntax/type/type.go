// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package _type

import (
	"reflect"

	"github.com/wuhan005/crud/internal/syntax"
)

type GoType interface {
	GoType() reflect.Type
}

type DefinedType interface {
	syntax.Syntax
	GoType
}

type BuiltInType interface {
	GoType
}

func GoTypeString(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem().String()
	}
	return typ.String()
}
