// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package variable

import (
	"bytes"
	"fmt"
	"reflect"
)

var _ Variable = (*Errors)(nil)

type Errors []*Error

func NewErrorsVar(errs ...*Error) Errors {
	return errs
}

func (e Errors) Name() string {
	return ""
}

func (e Errors) Decl() string {
	var buf bytes.Buffer
	buf.WriteString("var (")
	for _, e := range e {
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf("%s = errors.New(%q)", e.name, e.str))
	}
	buf.WriteString("\n)")
	return buf.String()
}

func (e Errors) DocString() string {
	return ""
}

func (e Errors) Body() string {
	return ""
}

func (e Errors) GoType() reflect.Type {
	return reflect.TypeOf(([]error)(nil))
}
