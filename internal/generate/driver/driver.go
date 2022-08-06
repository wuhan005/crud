// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package driver

import (
	"context"

	"github.com/wuhan005/crud/internal/db"
)

type Driver interface {
	GetStructure(ctx context.Context) (tables []string, tableColumnsSet map[string][]db.TableColumn, tableIndexesSet map[string][]db.TableIndex, _ error)
}
