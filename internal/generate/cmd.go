// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package generate

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/wuhan005/crud/internal/db"
	"github.com/wuhan005/crud/internal/db/backend"
	"github.com/wuhan005/crud/internal/dbutil"
)

type Generator interface {
	Generate(ctx context.Context) error
}

func Action(c *cli.Context) error {
	dsn := c.String("dsn")

	databaseType, err := dbutil.GetDatabaseType(dsn)
	if err != nil {
		return errors.Wrap(err, "get database type")
	}

	var databaseBackend db.Backend
	switch databaseType {
	case dbutil.TypePostgres:
		databaseBackend = backend.NewPostgresBackend(dsn)
	}

	// Connect to the database.
	if err := databaseBackend.Connect(c.Context); err != nil {
		return errors.Wrap(err, "connect to database")
	}
	defer func() { _ = databaseBackend.Close(c.Context) }()

	// Get all tables of the database.
	tables, err := databaseBackend.GetAllTables(c.Context)
	if err != nil {
		return errors.Wrap(err, "get all tables")
	}

	var tableColumnsSet map[string][]db.TableColumn
	var tableIndexesSet map[string][]db.TableIndex
	for _, tableName := range tables {
		tableColumns, err := databaseBackend.GetTableColumns(c.Context, tableName)
		if err != nil {
			return errors.Wrapf(err, "get table fields, table: %q", tableName)
		}
		tableColumnsSet[tableName] = tableColumns

		tableIndexes, err := databaseBackend.GetTableIndexes(c.Context, tableName)
		if err != nil {
			return errors.Wrapf(err, "get table indexes, table: %q", tableName)
		}
		tableIndexesSet[tableName] = tableIndexes
	}

	if err := Generate(c.Context, Options{
		Tables:          tables,
		TableIndexesSet: tableIndexesSet,
		TableColumnsSet: tableColumnsSet,
	}); err != nil {
		return errors.Wrap(err, "generate")
	}
	return nil
}
