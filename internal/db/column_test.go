// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColumnLower(t *testing.T) {
	for _, tc := range []struct {
		name   string
		column ColumnName
		want   ColumnName
	}{
		{
			name:   "upper",
			column: "EGGPLANT",
			want:   "eggplant",
		},
		{
			name:   "upper with underscore",
			column: "I_LOVE_eggpLANT",
			want:   "iLoveEggplant",
		},
		{
			name:   "upper with underscore and number",
			column: "I_LOVE_E99p1ant",
			want:   "iLoveE99P1Ant",
		},
		{
			name:   "camel case",
			column: "ILoveEggplant",
			want:   "iloveeggplant",
		},
		{
			name:   "special case",
			column: "user_id",
			want:   "userID",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.column.Lower()
			require.Equal(t, tc.want, got)
		})
	}
}

func TestColumnUpper(t *testing.T) {
	for _, tc := range []struct {
		name   string
		column ColumnName
		want   ColumnName
	}{
		{
			name:   "lower",
			column: "eggplant",
			want:   "Eggplant",
		},
		{
			name:   "lower with underscore",
			column: "I_LOVE_eggpLANT",
			want:   "ILoveEggplant",
		},
		{
			name:   "lower with underscore and number",
			column: "I_LOVE_E99p1ant",
			want:   "ILoveE99P1Ant",
		},
		{
			name:   "special case",
			column: "user_id",
			want:   "UserID",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.column.Upper()
			require.Equal(t, tc.want, got)
		})
	}
}
