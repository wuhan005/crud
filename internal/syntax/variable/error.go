// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package variable

import (
	"fmt"
	"reflect"
)

var _ Variable = (*Error)(nil)

type Error struct {
	name string
	str  string
}

func NewErrorVar(name, str string) *Error {
	return &Error{
		name: name,
		str:  str,
	}
}

func (e Error) Name() string {
	return e.name
}

func (e Error) Decl() string {
	return fmt.Sprintf("var %s = errors.New(%q)", e.name, e.str)
}

func (e Error) DocString() string {
	return fmt.Sprintf("// %s returns when %s", e.name, e.str)
}

func (e Error) GoType() reflect.Type {
	return reflect.TypeOf((*error)(nil))
}

func (e *Error) Body() string {
	return fmt.Sprintf("errors.New(%q)", e.str)
}
