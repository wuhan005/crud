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

var replaceMap = map[string]string{
	"id":  "ID",
	"uid": "UID",
}

type ColumnName string

// Lower transforms the lower snake case column name to lower camel case.
// Eggplant -> eggplant
// i_love_eggplant -> iLoveEggplant
// user_id -> userID
func (c ColumnName) Lower() ColumnName {
	name := c.replaceSpecialWord()
	return ColumnName(strcase.ToLowerCamel(name))
}

// Upper transforms the lower snake case column name to upper camel case.
// eggplant -> Eggplant
// user_id -> UserID
// userID -> UserID
func (c ColumnName) Upper() ColumnName {
	name := c.replaceSpecialWord()
	return ColumnName(strcase.ToCamel(name))
}

func (c ColumnName) String() string {
	return string(c)
}

func (c ColumnName) replaceSpecialWord() string {
	name := strings.ToLower(string(c))
	name = strcase.ToSnake(name)
	groups := strings.Split(name, "_")

	for i, group := range groups {
		if _, ok := replaceMap[group]; ok {
			groups[i] = replaceMap[group]
		}
	}

	return strings.Join(groups, "_")
}
