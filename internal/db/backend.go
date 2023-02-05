// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"reflect"
	"time"
)

// Backend is a database backend interface.
type Backend interface {
	// Connect connects to the database.
	Connect(ctx context.Context) error
	// Close closes the database connection.
	Close(ctx context.Context) error

	// GetAllTables returns all table names.
	GetAllTables(ctx context.Context) ([]string, error)
	// GetTableColumns returns all columns of the given table.
	GetTableColumns(ctx context.Context, tableName string) ([]TableColumn, error)
	// GetTableIndexes returns all indexes of the given table.
	GetTableIndexes(ctx context.Context, tableName string) ([]TableIndex, error)

	// ParseColumnType parses the given database column type to the defined type.
	ParseColumnType(typ string) (ColumnType, error)
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

// GetIndexColumns returns the index columns of the given table columns.
func (i TableIndex) GetIndexColumns(allColumns ...TableColumn) []TableColumn {
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
