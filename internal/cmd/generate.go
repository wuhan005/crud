// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/wuhan005/crud/internal/generate"
)

var Generate = &cli.Command{
	Name:    "generate",
	Aliases: []string{"gen"},
	Action:  generate.Action,
	Flags: []cli.Flag{
		stringFlag("dsn", "", "Database connection string"),
	},
}
