package main

import (
	"bytes"
	"fmt"
	"strings"

	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

// TypeField is struct
type TypeField struct {
	FieldName string
	FieldType string
	FieldTag  string
}

// StructVisitor is struct
type StructVisitor struct {
	name   string
	fields []TypeField
}

// Visit is func
func (v *StructVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}

	switch t := node.(type) {
	case *ast.TypeSpec:
		v.name = fmt.Sprintf("%s", t.Name)
	case *ast.StructType:
		fields := t.Fields.List
		v.fields = extractFields(fields)
	default:
	}

	return v
}

func extractFields(fields []*ast.Field) []TypeField {
	extractFields := []TypeField{}
	for _, f := range fields {
		typeField, ok := extractField(f)

		if ok {
			extractFields = append(extractFields, typeField)
		}
	}

	return extractFields
}

func extractField(field *ast.Field) (TypeField, bool) {
	if field.Tag == nil {
		return TypeField{}, false
	}

	// fieldName
	fieldName := field.Names[0].Name

	// fieldType
	var b bytes.Buffer
	format.Node(&b, token.NewFileSet(), field.Type)
	fieldType := b.String()

	// fieldTag
	var bb bytes.Buffer
	format.Node(&bb, token.NewFileSet(), field.Tag)
	fieldTag := bb.String()

	return TypeField{
		FieldName: fieldName,
		FieldType: fieldType,
		FieldTag:  fieldTag,
	}, true
}

func main() {
	source := `
		package Sample

		type Test struct {
			Id         Hoge.Fuga  'json:"id"'
			Name       string     'json:"name"'
			CreatedAt  *time.Time 'json:"createdAt"'
			UserStatus string
		}
	`

	source = strings.Replace(source, "'", "`", -1)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", source, 0)
	if err != nil {
		fmt.Println(err)
	}

	// ast.Print(nil, f)

	visitor := &StructVisitor{}
	ast.Walk(visitor, f)

	fmt.Println(visitor.name)
	fmt.Println(visitor.fields)
}
