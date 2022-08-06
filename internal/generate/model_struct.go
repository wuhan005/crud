// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package generate

import (
	"fmt"

	"github.com/wuhan005/crud/internal/db"
	"github.com/wuhan005/crud/internal/dbutil"
	_type "github.com/wuhan005/crud/internal/syntax/type"
)

type makeModelStructOptions struct {
	TableName dbutil.TableName
	Columns   []db.TableColumn
}

func makeModelStruct(opts makeModelStructOptions) (*_type.StructType, error) {
	tableName := opts.TableName
	structName := tableName.UpperSingular()
	docString := fmt.Sprintf("%s represents the %s.", structName, tableName.LowerPlural())

	structFields := make([]*_type.StructField, 0, len(opts.Columns))
	// Add gorm.Model as the first field.
	gormModelField := _type.NewStructField("", "gorm.Model", "")
	gormModelField.WithNewLine = true
	structFields = append(structFields, gormModelField)

	for _, column := range opts.Columns {
		columnName := column.Name
		if isIgnoreColumn(columnName) {
			continue
		}
		structFields = append(structFields, _type.NewStructField(columnName.Upper(), _type.GoTypeString(column.Type.GoType())))
	}

	structType := _type.NewStructType(_type.NewStructTypeOptions{
		Name:      structName,
		DocString: docString,
		Fields:    structFields,
	})
	return structType, nil
}

func isIgnoreColumn(columnName db.ColumnName) bool {
	s := struct{}{}
	_, ok := map[db.ColumnName]struct{}{
		"id":         s,
		"created_at": s,
		"updated_at": s,
		"deleted_at": s,
	}[columnName]
	return ok
}
