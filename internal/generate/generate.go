// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package generate

import (
	"bytes"
	"context"
	"text/template"

	"github.com/pkg/errors"

	"github.com/wuhan005/crud/internal/db"
	"github.com/wuhan005/crud/internal/syntax/function"
	_type "github.com/wuhan005/crud/internal/syntax/type"
)

const codeTemplate = `package {{.PackageName}}

import (
	"context"
	
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ {{.Users}}Store = (*{{.users}})(nil)

// {{.Users}} is the default instance of the {{.Users}}Store.
var {{.Users}} {{.Users}}Store

// {{.Users}}Store is the persistent interface for {{.users}}.
type {{.Users}}Store interface {
	{{.FunctionDecl}}
}

// New{{.Users}}Store returns a {{.Users}}Store instance with the given database connection.
func New{{.Users}}Store(db *gorm.DB) {{.Users}}Store {
	return &{{.users}}{db}
}

{{.ModelStruct}}

type {{.users}} struct {
	*gorm.DB
}

{{.FunctionBody}}
`

type generateTableCodeOptions struct {
	tableName db.TableName
	columns   []db.TableColumn
	indexes   []db.TableIndex
}

func generateTableCode(ctx context.Context, opts generateTableCodeOptions) ([]byte, error) {
	headerTemplate, err := template.New("header").Parse(codeTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parse header template")
	}

	// Make database model struct.
	model, err := makeModelStruct(makeModelStructOptions{
		TableName: opts.tableName,
		Columns:   opts.columns,
	})
	if err != nil {
		return nil, errors.Wrap(err, "make model struct")
	}

	functionGroups, err := makeFunctions(makeFunctionsOptions{
		TableName: opts.tableName,
		Model:     model,
		columns:   opts.columns,
		indexes:   opts.indexes,
	})
	if err != nil {
		return nil, errors.Wrap(err, "make functions")
	}
	var functionDecl, functionBody string
	for _, group := range functionGroups {
		// Add error declaration.
		for _, err := range group.Errors() {
			functionBody += err.Decl() + "\n"
		}
		for _, fun := range group.Functions() {
			functionDecl += fun.DocString() + "\n" + fun.Decl() + "\n"
			functionBody += fun.Body() + "\n\n"
		}
	}

	modelStruct := model.DocString() + "\n" + model.Decl()

	var headerData bytes.Buffer
	if err := headerTemplate.Execute(&headerData, map[string]interface{}{
		"PackageName":  "db",
		"ModelStruct":  modelStruct,
		"FunctionDecl": functionDecl,
		"FunctionBody": functionBody,
		"Users":        opts.tableName.UpperPlural(),
		"users":        opts.tableName.LowerPlural(),
		"User":         opts.tableName.UpperSingular(),
		"user":         opts.tableName.LowerSingular(),
	}); err != nil {
		return nil, errors.Wrap(err, "execute template")
	}
	return headerData.Bytes(), nil
}

type makeFunctionsOptions struct {
	TableName db.TableName
	Model     *_type.StructType
	columns   []db.TableColumn
	indexes   []db.TableIndex
}

func makeFunctions(options makeFunctionsOptions) ([]*function.Group, error) {
	functionGroups := make([]*function.Group, 0)

	// GetXXX
	getFunctionGroup := function.NewGroup()
	for _, index := range options.indexes {
		indexColumns := index.GetIndexColumns(options.columns...)
		getFunction, err := function.NewFunctionGet(function.NewFunctionGetOptions{
			Group:     getFunctionGroup,
			TableName: options.TableName,
			Model:     options.Model,
			Columns:   indexColumns,
			Indexes:   options.indexes,
		})
		if err != nil {
			return nil, errors.Wrap(err, "make get columns")
		}
		getFunctionGroup.AddFunction(getFunction)
	}
	functionGroups = append(functionGroups, getFunctionGroup)

	return functionGroups, nil
}
