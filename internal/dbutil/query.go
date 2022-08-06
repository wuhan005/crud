// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"reflect"

	"github.com/keegancsmith/sqlf"
)

func Expand(v interface{}) *sqlf.Query {
	kind := reflect.TypeOf(v).Kind()
	if kind != reflect.Slice {
		return sqlf.Sprintf("%s", v)
	}

	val := reflect.ValueOf(v)
	qs := make([]*sqlf.Query, val.Len())
	for i := 0; i < val.Len(); i++ {
		qs[i] = sqlf.Sprintf("%s", val.Index(i).Interface())
	}
	return sqlf.Join(qs, ",")
}
