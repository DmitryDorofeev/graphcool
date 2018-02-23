package templates

import "fmt"

func generateObjectMarshaler(typeName string) string {
	return fmt.Sprintf(`
			func (s *%sMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, *errors.QueryError) {
				buf := bytes.Buffer{}
				buf.WriteString("{")
				for i, value := range fields {
					field, ok := selections[i].(*query.Field)
					if !ok {
						continue
					}

					if (i != 0) {
						buf.WriteString(",")
					}

					buf.WriteString("\"" + field.Name.Name + "\"")
					buf.WriteString(":")
					buf.Write(value)
				}
				buf.WriteString("}")
				return buf.Bytes(), nil
			}
		`, typeName)
}

func generateListMarshaler(typeName string) string {
	return fmt.Sprintf(`
			func (s *%sMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, *errors.QueryError) {
				buf := bytes.Buffer{}
				buf.WriteString("[")
				for i, value := range fields {
					if (i != 0) {
						buf.WriteString(",")
					}

					buf.Write(value)
				}
				buf.WriteString("]")
				return buf.Bytes(), nil
			}
		`, typeName)
}
