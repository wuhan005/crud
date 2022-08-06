// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package function

import (
	"github.com/wuhan005/crud/internal/syntax/variable"
)

type Function interface {
	variable.Variable
	InputParameters() variable.Parameters
	OutputParameters() variable.Parameters
}
