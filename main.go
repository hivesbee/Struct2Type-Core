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

type TypeField struct {
	FieldName string
	FieldType string
	FieldTag  string
}

type StructVisitor struct {
	name   string
	fields []TypeField
}

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

	fmt.Println(extractFields)

	return extractFields
}

func extractField(field *ast.Field) (TypeField, bool) {
	if field.Tag == nil {
		return TypeField{}, false
	}

	// fieldName
	fieldName := field.Names[0].Name
	fmt.Println(fieldName)

	// fieldType
	var b bytes.Buffer
	format.Node(&b, token.NewFileSet(), field.Type)
	fieldType := b.String()
	fmt.Println(fieldType)

	// fieldTag
	var bb bytes.Buffer
	format.Node(&bb, token.NewFileSet(), field.Tag)
	fieldTag := bb.String()
	fmt.Println(fieldTag)

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
			Id        Hoge.Fuga     'json:"id"'
			Name      string     'json:"name"'
			CreatedAt *time.Time
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

	// ast.Inspect(f, func(n ast.Node) bool {
	// 	var s string
	// 	switch x := n.(type) {
	// 	case *ast.BasicLit:
	// 		s = x.Value
	// 	case *ast.Ident:
	// 		s = x.Name
	// 	}
	// 	if s != "" {
	// 		fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
	// 	}
	// 	return true
	// })
}
