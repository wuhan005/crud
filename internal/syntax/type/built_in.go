// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package _type

import (
	"context"
	"reflect"
)

type ContextType struct{}

func NewContextType() ContextType {
	return ContextType{}
}

func (t ContextType) GoType() reflect.Type {
	return reflect.TypeOf((*context.Context)(nil))
}

type ErrorType struct{}

func NewErrorType() ErrorType {
	return ErrorType{}
}

func (t ErrorType) GoType() reflect.Type {
	return reflect.TypeOf((*error)(nil))
}
