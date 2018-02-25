package graphcool

import "fmt"

type Field string
type Getter string
type Mutation string

type Argument struct {
	Name  string
	Value interface{}
}
type Arguments map[string]interface{}

func (a *Arguments) GetString(name string) (string, error) {
	arg, ok := (*a)[name]
	if !ok {
		return "", fmt.Errorf("Argument %s not present", name)
	}
	strArg, ok := arg.(string)
	if !ok {
		return "", fmt.Errorf("Cannot convert argument to string")
	}

	return strArg, nil
}
