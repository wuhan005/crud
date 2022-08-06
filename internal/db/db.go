// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"reflect"
	"time"
)

type DB interface {
	Conn(ctx context.Context) error
	GetAllTables(ctx context.Context) ([]string, error)
	GetTableColumn(ctx context.Context, tableName string) ([]TableColumn, error)
	GetTableIndexes(ctx context.Context, tableName string) ([]TableIndex, error)
	Close(ctx context.Context) error
}

type ColumnType string

const (
	ColumnTypeNumber  ColumnType = "number"
	ColumnTypeString  ColumnType = "string"
	ColumnTypeTime    ColumnType = "time"
	ColumnTypeBoolean ColumnType = "boolean"
)

func (t ColumnType) GoType() reflect.Type {
	switch t {
	case ColumnTypeNumber:
		return reflect.TypeOf(int64(0))
	case ColumnTypeString:
		return reflect.TypeOf("")
	case ColumnTypeTime:
		return reflect.TypeOf(time.Time{})
	case ColumnTypeBoolean:
		return reflect.TypeOf(true)
	}
	return nil
}

type TableColumn struct {
	TableName  string
	Name       ColumnName
	Type       ColumnType
	IsNullable bool
}

type TableIndex struct {
	TableName string
	Name      string
	IsUnique  bool
	Columns   []string
}

func (i TableIndex) GetColumns(allColumns ...TableColumn) []TableColumn {
	columns := make([]TableColumn, 0)

	for _, indexColumnName := range i.Columns {
		for _, column := range allColumns {
			column := column
			if indexColumnName == string(column.Name) {
				columns = append(columns, column)
				break
			}
		}
	}
	return columns
}
