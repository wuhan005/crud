// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package generate

import (
	"context"
	"go/format"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

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

	tableColumnsSet := make(map[string][]db.TableColumn)
	tableIndexesSet := make(map[string][]db.TableIndex)
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
		OutputPath: c.String("output"),

		Tables:          tables,
		TableIndexesSet: tableIndexesSet,
		TableColumnsSet: tableColumnsSet,
	}); err != nil {
		return errors.Wrap(err, "generate")
	}
	return nil
}

type Options struct {
	OutputPath string

	Tables          []string
	TableColumnsSet map[string][]db.TableColumn
	TableIndexesSet map[string][]db.TableIndex
}

func Generate(ctx context.Context, opts Options) error {
	for _, name := range opts.Tables {
		tableName := db.TableName(name)
		data, err := generateTableCode(ctx, generateTableCodeOptions{
			tableName: tableName,
			columns:   opts.TableColumnsSet[name],
			indexes:   opts.TableIndexesSet[name],
		})
		if err != nil {
			return errors.Wrapf(err, "generate table: %q", tableName)
		}

		formattedCode, err := format.Source(data)
		if err != nil {
			formattedCode = data
			log.Error("Failed to format generated code: %v", err)
		}

		fileName := tableName.SnakePlural() + ".go"
		outputFilePath := filepath.Join(opts.OutputPath, fileName)
		if err := os.WriteFile(outputFilePath, formattedCode, 0644); err != nil {
			return errors.Wrap(err, "write file")
		}
	}
	return nil
}
