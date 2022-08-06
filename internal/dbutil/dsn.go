// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"net/url"

	"github.com/pkg/errors"
)

func GetDatabaseType(dsn string) (Type, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", errors.Wrap(err, "url parse")
	}

	switch u.Scheme {
	case "postgres":
		return TypePostgres, nil
	default:
		return "", ErrUnexpectedDatabaseType
	}
}
