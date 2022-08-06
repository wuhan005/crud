// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package function

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/wuhan005/crud/internal/db"
	"github.com/wuhan005/crud/internal/dbutil"
	_type "github.com/wuhan005/crud/internal/syntax/type"
	"github.com/wuhan005/crud/internal/syntax/variable"
)

var (
	ErrEmptyColumn = errors.New("empty column")
	ErrNilModel    = errors.New("nil model")
)

var _ Function = (*Get)(nil)

type Get struct {
	tableName dbutil.TableName
	model     *_type.StructType
	columns   []db.TableColumn

	multiResult  bool
	hasOptions   bool
	optionStruct *_type.StructType
	errors       variable.Errors
}

type NewFunctionGetOptions struct {
	TableName dbutil.TableName
	Model     *_type.StructType
	Columns   []db.TableColumn
	Indexes   []db.TableIndex
}

func NewFunctionGet(options NewFunctionGetOptions) (*Get, error) {
	model := options.Model
	if model == nil {
		return nil, ErrNilModel
	}
	columns := options.Columns
	if len(options.Columns) == 0 {
		return nil, ErrEmptyColumn
	}
	tableName := options.TableName

	f := &Get{
		tableName: tableName,
		model:     model,
		columns:   columns,
	}
	f.errors = variable.Errors{f.notFoundError()}

	functionName := f.Name()

	// If the query parameter is more than two, use an option struct parameter to hold.
	if len(columns) > 2 {
		// TODO Check if the given columns are unique index.
		// If not, the GET function maybe return more than one result.

		f.hasOptions = true
		// Make option struct.
		optionStructName := functionName + "Options"
		docString := "" // TODO: add doc string.
		fields := make([]*_type.StructField, 0, len(columns))
		for _, column := range columns {
			fields = append(fields, _type.NewStructField(column.Name.Upper(), string(column.Type)))
		}

		f.optionStruct = _type.NewStructType(_type.NewStructTypeOptions{
			Name:      optionStructName,
			DocString: docString,
			Fields:    nil,
		})
	}

	return f, nil
}

func (f *Get) Name() string {
	name := "Get"

	if len(f.columns) == 1 {
		columnName := f.columns[0].Name
		name += "By" + columnName.Upper()
	} else {
		for _, column := range f.columns {
			name += column.Name.Upper()
		}
	}

	return name
}

func (f *Get) Decl() string {
	return fmt.Sprintf("%s%s %s", f.Name(), f.InputParameters(), f.OutputParameters())
}

func (f *Get) GoType() reflect.Type {
	return reflect.TypeOf(func() {})
}

func (f *Get) InputParameters() variable.Parameters {
	parameters := variable.Parameters{variable.ContextParameter()}
	for _, column := range f.columns {
		parameters = append(parameters, variable.NewParameterWithGoType(column.Name.Lower(), column.Type))
	}
	return parameters
}

func (f *Get) OutputParameters() variable.Parameters {
	parameters := variable.Parameters{}

	if f.multiResult {
		// TODO
	} else {
		parameters = append(parameters, variable.NewParameterWithTypeString("", f.model.Name(), true))
	}

	parameters = append(parameters, variable.ErrorParameter())
	return parameters
}

func (f *Get) DocString() string {
	var returnName, givenName string
	if f.multiResult {
		returnName = f.tableName.LowerPlural() + " list"
		givenName = "options"
	} else {
		returnName = "a " + f.tableName.LowerSingular()
		givenName = f.columns[0].Name.Lower()
	}

	str := fmt.Sprintf("// %s returns %s with the given %s.", f.Name(), returnName, givenName)
	str += "\n" + "// The zero value in the options will be ignored."
	return str
}

func (f *Get) Body() string {
	errorDecl := f.errors.Decl()
	body := f.makeBody()
	str := errorDecl + "\n"
	str += fmt.Sprintf("func (db *%s) %s {\n %s \n}", f.tableName.LowerPlural(), f.Decl(), body)
	return str
}

func (f *Get) makeBody() string {
	modelName := f.tableName.UpperSingular()
	var resultVar, whereExpr, getExpr string
	if f.hasOptions {
		// TODO
	} else {
		column := f.columns[0]
		resultVar = f.tableName.LowerSingular()
		whereExpr = fmt.Sprintf(`"%s = ?", %s`, column.Name.Lower(), column.Name.Lower())
		getExpr = fmt.Sprintf("First(&%s)", resultVar)
	}

	str := fmt.Sprintf("var %s %s\n", resultVar, f.tableName.UpperSingular())
	str += fmt.Sprintf(`if err := db.WithContext(ctx).Model(&%s{}).Where(%s).%s.Error; err != nil {`, modelName, whereExpr, getExpr)
	str += "if errors.Is(err, gorm.ErrRecordNotFound) {\n"
	str += fmt.Sprintf("return nil, %s\n", f.notFoundError().Name())
	str += "}\n"
	str += "}\n"
	str += fmt.Sprintf(`return &%s, nil`, resultVar)
	return str
}

func (f *Get) notFoundError() *variable.Error {
	name := fmt.Sprintf("Err%sNotExists", f.tableName.UpperSingular())
	str := fmt.Sprintf("%s dose not exist", f.tableName.LowerSingular())
	return variable.NewErrorVar(name, str)
}
