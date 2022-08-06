// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package generate

import (
	"bytes"
	"context"
	"go/format"
	"os"
	"text/template"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/crud/internal/db"
	"github.com/wuhan005/crud/internal/dbutil"
	"github.com/wuhan005/crud/internal/syntax/function"
	_type "github.com/wuhan005/crud/internal/syntax/type"
)

const header = `package {{.PackageName}}

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

{{.ModelStructDocString}}
{{.ModelStructDecl}}

type {{.users}} struct {
	*gorm.DB
}

{{.FunctionBody}}
`

type Options struct {
	Tables          []string
	TableColumnsSet map[string][]db.TableColumn
	TableIndexesSet map[string][]db.TableIndex
}

func Generate(ctx context.Context, opts Options) error {
	for _, name := range opts.Tables {
		tableName := dbutil.TableName(name)
		data, err := generateTable(ctx, generateTableOptions{
			tableName: tableName,
			columns:   opts.TableColumnsSet[name],
			indexes:   opts.TableIndexesSet[name],
		})
		if err != nil {
			return errors.Wrap(err, "generate table")
		}

		formatCode, err := format.Source(data)
		if err != nil {
			formatCode = data
			log.Error("Failed to format generated code: %v", err)
		}

		outputFilePath := tableName.LowerPlural() + ".go"
		if err := os.WriteFile(outputFilePath, formatCode, 0644); err != nil {
			return errors.Wrap(err, "write file")
		}
	}
	return nil
}

type generateTableOptions struct {
	tableName dbutil.TableName
	columns   []db.TableColumn
	indexes   []db.TableIndex
}

func generateTable(ctx context.Context, opts generateTableOptions) ([]byte, error) {
	headerTemplate, err := template.New("header").Parse(header)
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

	functions, err := makeFunctions(makeFunctionsOptions{
		TableName: opts.tableName,
		Model:     model,
		columns:   opts.columns,
		indexes:   opts.indexes,
	})
	if err != nil {
		return nil, errors.Wrap(err, "make functions")
	}
	var functionDecl, functionBody string
	for _, fun := range functions {
		functionDecl += fun.DocString() + "\n" + fun.Decl() + "\n"
		functionBody += fun.Body() + "\n"
	}

	var headerData bytes.Buffer
	if err := headerTemplate.Execute(&headerData, map[string]interface{}{
		"PackageName":          "db",
		"FunctionDecl":         functionDecl,
		"ModelStructDocString": model.DocString(),
		"ModelStructDecl":      model.Decl(),
		"FunctionBody":         functionBody,
		"Users":                opts.tableName.UpperPlural(),
		"users":                opts.tableName.LowerPlural(),
		"User":                 opts.tableName.UpperSingular(),
		"user":                 opts.tableName.LowerSingular(),
	}); err != nil {
		return nil, errors.Wrap(err, "execute template")
	}
	return headerData.Bytes(), nil
}

type makeFunctionsOptions struct {
	TableName dbutil.TableName
	Model     *_type.StructType
	columns   []db.TableColumn
	indexes   []db.TableIndex
}

func makeFunctions(options makeFunctionsOptions) ([]function.Function, error) {
	functions := make([]function.Function, 0)

	// GetXXX
	for _, index := range options.indexes {
		indexColumns := index.GetColumns(options.columns...)
		getFunction, err := function.NewFunctionGet(function.NewFunctionGetOptions{
			TableName: options.TableName,
			Model:     options.Model,
			Columns:   indexColumns,
			Indexes:   options.indexes,
		})
		if err != nil {
			return nil, errors.Wrap(err, "make get columns")
		}
		functions = append(functions, getFunction)
	}

	return functions, nil
}
