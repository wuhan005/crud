// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package _type

import (
	"bytes"
	"fmt"
	"reflect"
)

var _ DefinedType = (*StructType)(nil)

type StructField struct {
	Name        string
	Type        string
	Tag         string
	WithNewLine bool
}

func NewStructField(name, typ string, tags ...string) *StructField {
	var tag string
	if len(tags) > 0 {
		tag = tags[0]
	}

	return &StructField{
		Name: name,
		Type: typ,
		Tag:  tag,
	}
}

func (s *StructField) String() string {
	tag := s.Tag
	if tag != "" {
		tag = "`" + tag + "`"
	}
	str := s.Name + " " + s.Type + " " + tag

	if s.WithNewLine {
		str += "\n"
	}
	return str
}

type StructType struct {
	name      string
	docString string
	fields    []*StructField
}

type NewStructTypeOptions struct {
	Name      string
	DocString string
	Fields    []*StructField
}

func NewStructType(opts NewStructTypeOptions) *StructType {
	return &StructType{
		name:      opts.Name,
		docString: opts.DocString,
		fields:    opts.Fields,
	}
}

func (s *StructType) Name() string {
	return s.name
}

func (s *StructType) Decl() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("type %s struct {\n", s.name))
	for _, field := range s.fields {
		buf.WriteString("	" + field.String())
		buf.WriteString("\n")
	}
	buf.WriteString("}")
	return buf.String()
}

func (s *StructType) DocString() string {
	return "// " + s.docString
}

func (s *StructType) GoType() reflect.Type {
	return reflect.TypeOf((*struct{})(nil))
}
