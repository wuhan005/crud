// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package variable

import (
	_type "github.com/wuhan005/crud/internal/syntax/type"
)

const EmptyParameterPlaceholder = "_"

type Parameters []*Parameter

func (p Parameters) String() string {
	str := "("
	for i, parameter := range p {
		str += parameter.String()
		if i != len(p)-1 {
			str += ", "
		}
	}
	str += ")"
	return str
}

// Parameter is a function parameter.
type Parameter struct {
	Name       string
	TypeString string
	Pointer    bool
}

func NewParameterWithGoType(name string, typ _type.GoType, pointer ...bool) *Parameter {
	return &Parameter{
		Name:       name,
		TypeString: _type.GoTypeString(typ.GoType()),
		Pointer:    len(pointer) > 0 && pointer[0],
	}
}

func NewParameterWithTypeString(name string, typeString string, pointer ...bool) *Parameter {
	return &Parameter{
		Name:       name,
		TypeString: typeString,
		Pointer:    len(pointer) > 0 && pointer[0],
	}
}

func (p Parameter) String() string {
	pointerStar := ""
	if p.Pointer {
		pointerStar = "*"
	}
	return pointerStar + p.Name + " " + p.TypeString
}

func ContextParameter() *Parameter {
	return NewParameterWithGoType("ctx", _type.NewContextType())
}

func ErrorParameter(name ...string) *Parameter {
	if len(name) == 0 {
		return NewParameterWithGoType("", _type.NewErrorType())
	}
	return NewParameterWithGoType(name[0], _type.NewErrorType())
}
