// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package variable

import (
	"reflect"

	"github.com/wuhan005/crud/internal/syntax"
)

type Variable interface {
	syntax.Syntax
	// Body is the body of the variable, such as function body, struct body, variable value, etc.
	Body() string
	GoType() reflect.Type
}
