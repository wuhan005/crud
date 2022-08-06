// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package generate

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/wuhan005/crud/internal/dbutil"
	"github.com/wuhan005/crud/internal/generate/driver"
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

	var structDriver driver.Driver
	switch databaseType {
	case dbutil.TypePostgres:
		structDriver = driver.NewPostgresDriver(dsn)
	}

	tables, tableColumnsSet, tableIndexesSet, err := structDriver.GetStructure(c.Context)
	if err != nil {
		return errors.Wrap(err, "get structure")
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
