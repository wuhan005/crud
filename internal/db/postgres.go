// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/keegancsmith/sqlf"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

type Postgres struct {
	DSN string

	db *sqlx.DB
}

func (p *Postgres) Conn(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", p.DSN)
	if err != nil {
		return errors.Wrap(err, "connect")
	}

	p.db = db
	return nil
}

func (p *Postgres) GetAllTables(ctx context.Context) ([]string, error) {
	var tables []string
	if err := p.db.SelectContext(ctx, &tables, `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'`); err != nil {
		return nil, errors.Wrap(err, "select")
	}
	return tables, nil
}

func (p *Postgres) GetTableColumn(ctx context.Context, tableName string) ([]TableColumn, error) {
	type postgresColumn struct {
		TableName  string `db:"table_name"`
		ColumnName string `db:"column_name"`
		DataType   string `db:"data_type"`
		IsNullable string `db:"is_nullable"`
	}

	q := sqlf.Sprintf(`SELECT table_name, column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = %s ORDER BY ordinal_position`, tableName)
	var columns []postgresColumn
	if err := p.db.SelectContext(ctx, &columns, q.Query(sqlf.PostgresBindVar), q.Args()...); err != nil {
		return nil, errors.Wrap(err, "select")
	}

	tableColumns := make([]TableColumn, 0, len(columns))
	for _, column := range columns {
		dataType, err := ParseColumnType(column.DataType)
		if err != nil {
			log.Warn("Unexpected column type: %q in table %q, ignore.", column.DataType, tableName)
			continue
		}
		columnName := ColumnName(column.ColumnName)

		tableColumns = append(tableColumns, TableColumn{
			TableName:  column.TableName,
			Name:       columnName,
			Type:       dataType,
			IsNullable: column.IsNullable == "YES",
		})
	}
	return tableColumns, nil
}

func (p *Postgres) GetTableIndexes(ctx context.Context, tableName string) ([]TableIndex, error) {
	type postgresIndex struct {
		TableName string `db:"tablename"`
		IndexName string `db:"indexname"`
		IndexDef  string `db:"indexdef"`
	}

	q := sqlf.Sprintf(`SELECT tablename, indexname, indexdef FROM pg_indexes WHERE schemaname = CURRENT_SCHEMA() AND tablename = %s;`, tableName)
	var indexes []postgresIndex
	if err := p.db.SelectContext(ctx, &indexes, q.Query(sqlf.PostgresBindVar), q.Args()...); err != nil {
		return nil, errors.Wrap(err, "select")
	}

	tableIndexes := make([]TableIndex, 0, len(indexes))
	for _, index := range indexes {
		indexDef := strings.ToLower(strings.TrimSpace(index.IndexDef))
		// Parse the index definition.
		sqlFields := strings.Fields(indexDef)

		if sqlFields[0] != "create" {
			return nil, errors.Errorf("unexpected index def: %q", indexDef)
		}

		isUnique := sqlFields[1] == "unique" && sqlFields[2] == "index"

		columnsStr := sqlFields[len(sqlFields)-1]
		columnsStr = strings.Trim(columnsStr, "()")
		columns := strings.Split(columnsStr, ",")

		tableIndexes = append(tableIndexes, TableIndex{
			TableName: index.TableName,
			Name:      index.IndexName,
			IsUnique:  isUnique,
			Columns:   columns,
		})
	}
	return tableIndexes, nil
}

func (p *Postgres) Close(_ context.Context) error {
	return p.db.Close()
}
