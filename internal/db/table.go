// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type TableName string

func (t TableName) LowerSingular() string {
	name := strings.ToLower(string(t))
	singular := inflection.Singular(name)
	return strcase.ToLowerCamel(singular)
}

func (t TableName) LowerPlural() string {
	name := strings.ToLower(string(t))
	singular := inflection.Plural(name)
	return strcase.ToLowerCamel(singular)
}

func (t TableName) UpperSingular() string {
	name := strings.ToLower(string(t))
	singular := inflection.Singular(name)
	return strcase.ToCamel(singular)
}

func (t TableName) UpperPlural() string {
	name := strings.ToLower(string(t))
	singular := inflection.Plural(name)
	return strcase.ToCamel(singular)
}
