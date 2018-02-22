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
		func (s *%sMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, *errors.QueryError) {
			fmt.Println("Looking up %s")
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
		func (s *%sMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, *errors.QueryError) {
			fmt.Println("Looking up %s")
			res := make([][]byte, 0)
			if len(selections) == 0 {
				return nil, errors.Errorf("Objects must have selections (field %s has no selections)")
			}
			for _, item := range s.Value {
				f := %sMeta{
					Value: item,
				}

				b, _ := f.Lookup(ctx, selections)
				res = append(res, b)
			}
			return s.Marshal(ctx, selections, res)
		}
	`, typeName, typeName, typeName, itemType)
}

// func generateListCases(fields []parser.StructField, s parser.ParsedStructs) (cases string) {
// 	for _, field := range fields {
// 		t := field.Type()
// 		tag := field.Tag().Get("graphql")
// 		fieldName, _ := getFieldName(tag)
// 		switch t.(type) {
// 		case *types.Named:
// 			cases += fmt.Sprintf(`
// 					case "%s":
// 				`, fieldName)
// 		case *types.Basic:
// 			cases += fmt.Sprintf(`
// 					case "%s":
// 						fmt.Println("touch %s")
// 						f := %sMeta{
// 							Value: s.Value.%s,
// 						}

// 						res = append(res, f.Marshal())
// 				`, fieldName)
// 		}
// 	}
// 	return
// }

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

func generateCases(fields []parser.StructField, s parser.ParsedStructs) (lookups string) {
	for _, field := range fields {
		t := field.Type()
		n := field.Name()

		tag := field.Tag().Get("graphql")
		fieldName, typeName := getFieldName(tag)

		switch t.(type) {
		case *types.Named:
			if typeName == "" {
				typeName = t.(*types.Named).Obj().Name()
			}

			lookups += fmt.Sprintf(`
				case "%s":
					f := %sMeta{
					}

					err := f.Value.Resolve(ctx)
					if err != nil {
						return nil, err
					}

					innerField, err := f.Lookup(ctx, field.Selections)
					if err != nil {
						return nil, err
					}

					res = append(res, innerField)

			`, fieldName, typeName)
		case *types.Basic:
			if typeName == "" {
				typeName = strings.Title(t.(*types.Basic).Name())
			}

			lookups += fmt.Sprintf(`
				case "%s":
					fmt.Println("touch %s")
					f := %sMeta{
						Value: s.Value.%s,
					}

					res = append(res, f.Marshal())
			`, fieldName, fieldName, typeName, n)
		}

	}
	return
}
