package parser

import (
	"strings"
	"text/scanner"

	"github.com/DmitryDorofeev/graphcool/common"
)

type Field struct {
	Name       string
	Args       common.InputValueList
	Type       common.Type
	Directives common.DirectiveList
	Desc       string
}

type FieldList []*Field

func ParseField(tag string) *Field {
	sc := &scanner.Scanner{
		Mode: scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings,
	}
	sc.Init(strings.NewReader(tag))

	l := common.New(sc)

	f := &Field{}
	f.Desc = l.DescComment()
	f.Name = l.ConsumeIdent()
	if l.Peek() == '(' {
		l.ConsumeToken('(')
		for l.Peek() != ')' {
			f.Args = append(f.Args, common.ParseInputValue(l))
		}
		l.ConsumeToken(')')
	}
	l.ConsumeToken(':')
	f.Type = common.ParseType(l)
	f.Directives = common.ParseDirectives(l)
	return f
}
