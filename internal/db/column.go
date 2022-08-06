// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

var ErrUnexpectedColumnType = errors.New("unexpected column type")

func ParseColumnType(typ string) (ColumnType, error) {
	v, ok := map[string]ColumnType{
		"character":              ColumnTypeString,
		"integer":                ColumnTypeNumber,
		"time without time zone": ColumnTypeTime,
	}[typ]
	if !ok {
		return "", ErrUnexpectedColumnType
	}
	return v, nil
}

type ColumnName string

func (c *ColumnName) Lower() string {
	name := strings.ToLower(string(*c))
	return strcase.ToLowerCamel(name)
}

func (c *ColumnName) Upper() string {
	name := strings.ToLower(string(*c))

	replaceMap := map[string]string{
		"id":  "ID",
		"uid": "UID",
	}
	var replaced bool

	groups := strings.Split(strcase.ToSnake(name), "_")
	for i, group := range groups {
		if _, ok := replaceMap[group]; ok {
			groups[i] = replaceMap[group]
			replaced = true
		}
	}
	// If there is only one element, and it has been replaced, return it directly.
	// Otherwise, the first element will be transformed to camel case.
	if len(groups) == 1 && replaced {
		return groups[0]
	}

	name = strings.Join(groups, "_")
	return strcase.ToCamel(name)
}
