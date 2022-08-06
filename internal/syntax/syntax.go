// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package syntax

// Syntax is the basic element in the Go source code.
// Such as a function, a struct, a parameter, a variable, etc.
type Syntax interface {
	// Name is the name of the syntax.
	Name() string
	// Decl is the declaration of the syntax.
	Decl() string
	// DocString is the documentation string comments of the syntax.
	DocString() string
}
