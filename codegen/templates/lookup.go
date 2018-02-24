package templates

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/DmitryDorofeev/graphcool/common"
	"github.com/DmitryDorofeev/graphcool/parser"
)

func generateLookup(typeName, cases string) string {
	return fmt.Sprintf(`
		func (s *%sMeta) Lookup(ctx context.Context, selections []query.Selection, vars map[string] interface{}) ([]byte, *errors.QueryError) {
			if len(selections) == 0 {
				return nil, errors.Errorf("Objects must have selections (field %s has no selections)")
			}
			res := make([][]byte, 0)
			for _, selection := range selections {
				field, ok := selection.(*query.Field)
				if !ok {
					fmt.Println("cannot cast to Field")
					continue
				}
				switch field.Name.Name {
					%s
					default:
						return nil, errors.Errorf("unknown field " + field.Name.Name)
				}
			}
			return s.Marshal(ctx, selections, res)
		}
	`, typeName, typeName, cases)
}

func generateListLookup(typeName, itemType string) string {
	return fmt.Sprintf(`
		func (s *%sMeta) Lookup(ctx context.Context, selections []query.Selection, vars map[string]interface{}) ([]byte, *errors.QueryError) {
			if len(selections) == 0 {
				return nil, errors.Errorf("Objects must have selections (field %s has no selections)")
			}
			res := make([][]byte, 0)
			for _, item := range s.Value {
				f := %sMeta{
					Value: item,
				}

				b, _ := f.Lookup(ctx, selections, vars)
				res = append(res, b)
			}
			return s.Marshal(ctx, selections, res)
		}
	`, typeName, typeName, itemType)
}

func getFieldName(tag string) (fieldName string, typeName string) {
	if strings.Contains(tag, ":") {
		fieldInfo := parser.ParseField(tag)
		fieldName = fieldInfo.Name
		switch fieldInfo.Type.(type) {
		case *common.TypeName:
			typeName = fieldInfo.Type.(*common.TypeName).Name
		case *common.List:
			typeName = fieldInfo.Type.(*common.List).OfType.(*common.TypeName).Name
		}
	} else {
		fieldName = tag
	}
	return
}

func rangeFields(cur parser.ParsedStruct, s parser.ParsedStructs) (fields string) {
	for _, method := range cur.Methods {
		if method.Doc == nil {
			continue
		}

		fieldName, typeName := getFieldName(strings.TrimSpace(strings.TrimPrefix(method.Doc.List[0].Text, "//")))
		fields += generateMethod(fieldName, typeName, method.Name)
	}

	for _, field := range cur.Fields {
		t := field.Type()
		n := field.Name()

		tag := field.Tag().Get("graphql")
		fieldName, typeName := getFieldName(tag)
		fields += generateField(t, n, typeName, fieldName)
	}
	return
}

func generateMethod(fieldName, typeName, methodName string) string {
	return generateComplexField(fieldName, typeName, methodName)
}

func generateField(t types.Type, n, typeName, fieldName string) (cases string) {
	switch t.(type) {
	case *types.Named:
		if typeName == "" {
			typeName = t.(*types.Named).Obj().Name()
		}

		cases += generateComplexField(fieldName, typeName, "")
	case *types.Basic:
		if typeName == "" {
			typeName = strings.Title(t.(*types.Basic).Name())
		}

		cases += generateScalarField(fieldName, typeName, n)
	}
	return
}

func generateComplexField(fieldName, typeName, methodName string) string {
	exec := "err := f.Value.Resolve(ctx, s.Value, args)"
	if methodName != "" {
		exec = fmt.Sprintf("val, err := s.Value.%s(ctx, s.Value, args)\nf.Value = val", methodName)
	}
	return fmt.Sprintf(`
				case "%s":
					f := %sMeta{
					}

					args := make(graphql.Arguments, 0)
					for _, arg := range field.Arguments {
						args[arg.Name.Name] = arg.Value.Value(vars)
					}

					%s
					if err != nil {
						return nil, err
					}

					innerField, err := f.Lookup(ctx, field.Selections, vars)
					if err != nil {
						return nil, err
					}

					res = append(res, innerField)

			`, fieldName, typeName, exec)
}

func generateScalarField(fieldName, typeName, structFieldName string) string {
	return fmt.Sprintf(`
		case "%s":
			f := %sMeta{
				Value: s.Value.%s,
			}

			res = append(res, f.Marshal())
	`, fieldName, typeName, structFieldName)
}
