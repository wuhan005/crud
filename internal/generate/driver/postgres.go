// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package driver

import (
	"context"

	"github.com/pkg/errors"

	"github.com/wuhan005/crud/internal/db"
)

type postgres struct {
	*db.Postgres
}

func NewPostgresDriver(dsn string) *postgres {
	return &postgres{
		Postgres: &db.Postgres{
			DSN: dsn,
		},
	}
}

func (p *postgres) GetStructure(ctx context.Context) (tables []string, tableColumnsSet map[string][]db.TableColumn, tableIndexesSet map[string][]db.TableIndex, _ error) {
	if err := p.Conn(ctx); err != nil {
		return nil, nil, nil, errors.Wrap(err, "connect")
	}
	defer func() { _ = p.Close(ctx) }()

	tables, err := p.GetAllTables(ctx)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "get tables")
	}

	tableColumnsSet = make(map[string][]db.TableColumn, len(tables))
	tableIndexesSet = make(map[string][]db.TableIndex, len(tables))

	for _, tableName := range tables {
		tableColumns, err := p.GetTableColumn(ctx, tableName)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "get table fields, table: %q", tableName)
		}
		tableColumnsSet[tableName] = tableColumns

		tableIndexes, err := p.GetTableIndexes(ctx, tableName)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "get table indexes, table: %q", tableName)
		}
		tableIndexesSet[tableName] = tableIndexes
	}

	return tables, tableColumnsSet, tableIndexesSet, nil
}
