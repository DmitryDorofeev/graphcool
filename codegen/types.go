package codegen

import (
	"fmt"
	"go/ast"
	"reflect"
)

type Schema struct {
	Pkg   string
	Types map[string]Type
	Tree  map[string]Element
}

type Field struct {
	Name string
	Type string
}

type Type struct {
	Name   string
	Fields []Field
}

type Element struct {
	Name     string
	Type     string
	Children []Element
}

func extractType(e ast.Expr) (string, string) {
	switch rt := e.(type) {
	case *ast.SelectorExpr:
		full, _ := extractType(rt.X)
		return full + "." + rt.Sel.Name, rt.Sel.Name
	case *ast.Ident:
		return rt.Name, rt.Name
	case *ast.StarExpr:
		full, ident := extractType(rt.X)
		return "*" + full, ident
	default:
		return "", ""
	}
}

func ExtractField(field *ast.Field) *Field {
	if field.Tag == nil {
		return nil
	}

	tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
	_, shortType := extractType(field.Type)
	gqlTag := tag.Get("graphql")

	if gqlTag == "" {
		return nil
	}

	return &Field{
		Name: tag.Get("graphql"),
		Type: shortType,
	}
}

func ExtractStruct(spec ast.Spec) (*Type, error) {
	currType, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("%#T is not ast.TypeSpec\n", spec)
	}

	currStruct, ok := currType.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("%#T is not ast.StructType\n", currStruct)
	}

	typeName := currType.Name.Name
	if typeName == "" {
		return nil, fmt.Errorf("invalid typename")
	}

	t := &Type{
		Name: typeName,
	}

	for _, field := range currStruct.Fields.List {
		f := ExtractField(field)
		if f != nil {
			t.Fields = append(t.Fields, *f)
		}
	}

	if len(t.Fields) == 0 {
		return nil, fmt.Errorf("no gql fields in struct")
	}

	return t, nil
}

func ExtractGQLStructs(decl ast.Decl) (map[string]Type, error) {
	types := make(map[string]Type)
	g, ok := decl.(*ast.GenDecl)
	if !ok {
		return types, fmt.Errorf("%#v is not *ast.GenDecl", decl)
	}

	for _, spec := range g.Specs {
		t, err := ExtractStruct(spec)
		if err != nil {
			continue
		}
		types[t.Name] = *t
	}

	return types, nil
}
