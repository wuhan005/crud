// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package function

import (
	"github.com/wuhan005/crud/internal/syntax/variable"
)

type Group struct {
	errors    []*variable.Error
	functions []Function
}

func NewGroup() *Group {
	return &Group{}
}

func (g *Group) Errors() []*variable.Error {
	return g.errors
}

func (g *Group) AddErrors(err *variable.Error) {
	for _, e := range g.errors {
		if e.Body() == err.Body() {
			return
		}
	}
	g.errors = append(g.errors, err)
}

func (g *Group) Functions() []Function {
	return g.functions
}

func (g *Group) AddFunction(f Function) {
	g.functions = append(g.functions, f)
}
