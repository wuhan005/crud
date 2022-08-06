// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"github.com/pkg/errors"
)

var ErrUnexpectedDatabaseType = errors.New("unexpected database type")

type Type string

const (
	TypePostgres Type = "postgres"
)
