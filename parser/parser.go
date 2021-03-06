package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/tools/go/loader"
)

// StructField represents one field in struct
type StructField struct {
	name string            //
	typ  types.Type        // field/method/parameter type
	tag  reflect.StructTag // field tag; or nil
}

type Method struct {
	Name string
	Doc  *ast.CommentGroup
}

func (sf StructField) Name() string {
	return sf.name
}

func (sf StructField) Type() types.Type {
	return sf.typ
}

func (sf StructField) Tag() reflect.StructTag {
	return sf.tag
}

// ParsedStructs is a map from struct type name to list of fields
type ParsedStructs map[string]ParsedStruct

// ParsedStruct represents struct info
type ParsedStruct struct {
	TypeName string
	Fields   []StructField
	Doc      *ast.CommentGroup // line comments; or nil
	Methods  []Method
}

func fileNameToPkgName(filePath, absFilePath string) string {
	dir := filepath.Dir(absFilePath)
	gopath := os.Getenv("GOPATH")
	if !strings.HasPrefix(dir, gopath) {
		// not in GOPATH
		return "./" + filepath.Dir(filePath)
	}

	r := strings.TrimPrefix(dir, gopath)
	r = strings.TrimPrefix(r, "/")  // may be and may not be
	r = strings.TrimPrefix(r, "\\") // may be and may not be
	r = strings.TrimPrefix(r, "src/")
	r = strings.TrimPrefix(r, "src\\")
	return r
}

func typeCheckFuncBodies(path string) bool {
	return false // don't type-check func bodies to speedup parsing
}

func loadProgramFromPackage(pkgFullName string) (*loader.Program, error) {
	// The loader loads a complete Go program from source code.
	conf := loader.Config{
		ParserMode:          parser.ParseComments,
		TypeCheckFuncBodies: typeCheckFuncBodies,
	}
	conf.Import(pkgFullName)
	lprog, err := conf.Load()
	if err != nil {
		return nil, fmt.Errorf("can't load program from package %q: %s",
			pkgFullName, err)
	}

	return lprog, nil
}

type structNamesInfo map[string]*ast.GenDecl
type methodsInfo map[string]map[string]*ast.FuncDecl

type structNamesVisitor struct {
	names      structNamesInfo
	methods    methodsInfo
	curGenDecl *ast.GenDecl
}

func (v *structNamesVisitor) Visit(n ast.Node) (w ast.Visitor) {
	switch n := n.(type) {
	case *ast.GenDecl:
		v.curGenDecl = n
	case *ast.TypeSpec:
		if _, ok := n.Type.(*ast.StructType); ok {
			v.names[n.Name.Name] = v.curGenDecl
		}

		if i, ok := n.Type.(*ast.Ident); ok {
			typ, ok := v.names[i.Name]
			if !ok {
				return v
			}
			v.names[n.Name.Name] = typ
		}

		if i, ok := n.Type.(*ast.ArrayType); ok {
			typ, ok := v.names[i.Elt.(*ast.Ident).Name]
			if !ok {
				return v
			}

			v.names[n.Name.Name] = typ
		}

	case *ast.FuncDecl:
		if n.Recv == nil {
			return v
		}

		// so dangerous
		receiverName := n.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
		if v.methods[receiverName] == nil {
			v.methods[receiverName] = make(map[string]*ast.FuncDecl)
		}

		v.methods[receiverName][n.Name.Name] = n
	}

	return v
}

func getStructInfo(fname string) (*structNamesVisitor, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("can't parse file %q: %s", fname, err)
	}

	v := &structNamesVisitor{
		names:   structNamesInfo{},
		methods: methodsInfo{},
	}
	ast.Walk(v, f)
	return v, nil
}

// GetStructsInFile lists all structures in file passed and returns them with all fields
func GetStructsInFile(filePath string) (*loader.PackageInfo, ParsedStructs, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("can't get abs path for %s", filePath)
	}

	neededStructs, err := getStructInfo(absFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("can't get struct names: %s", err)
	}

	packageFullName := fileNameToPkgName(filePath, absFilePath)
	lprog, err := loadProgramFromPackage(packageFullName)
	if err != nil {
		return nil, nil, err
	}

	pkgInfo := lprog.Package(packageFullName)
	if pkgInfo == nil {
		return nil, nil, fmt.Errorf("can't load types for file %s in package %q",
			filePath, packageFullName)
	}

	ret := ParsedStructs{}

	scope := pkgInfo.Pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)

		if neededStructs.names[name] == nil {
			continue
		}

		t := obj.Type().(*types.Named)

		if s, ok := t.Underlying().(*types.Struct); ok {
			parsedStruct := parseStruct(s, neededStructs, name)
			if parsedStruct != nil {
				parsedStruct.TypeName = name
				ret[name] = *parsedStruct
			}
		}

		if a, ok := t.Underlying().(*types.Slice); ok {
			realName := a.Elem().(*types.Named).Obj().Name()
			parsedStruct := ParsedStruct{
				TypeName: realName,
			}
			ret[name] = parsedStruct
		}

	}

	return pkgInfo, ret, nil
}

func newStructField(f *types.Var, tag string) *StructField {
	return &StructField{
		name: f.Name(),
		typ:  f.Type(),
		tag:  reflect.StructTag(tag),
	}
}

func parseStructFields(s *types.Struct) []StructField {
	var fields []StructField
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		if _, ok := f.Type().Underlying().(*types.Interface); ok {
			// skip interfaces
			continue
		}

		t := s.Tag(i)
		if reflect.StructTag(t).Get("graphql") == "" {
			continue
		}

		if f.Anonymous() {
			e, ok := f.Type().Underlying().(*types.Struct)
			if !ok {
				continue
			}

			pf := parseStructFields(e)
			if len(pf) == 0 {
				continue
			}

			fields = append(fields, pf...)
			continue
		}

		if !f.Exported() {
			continue
		}

		sf := newStructField(f, t)
		fields = append(fields, *sf)
	}

	return fields
}

func parseStruct(s *types.Struct, v *structNamesVisitor, name string) *ParsedStruct {
	fields := parseStructFields(s)

	var doc *ast.CommentGroup
	if decl := v.names[name]; decl != nil { // decl can be nil for embedded structs
		doc = decl.Doc // can obtain doc only from AST
	}

	methods := make([]Method, 0)

	for _, m := range v.methods[name] {
		methods = append(methods, Method{
			Name: m.Name.Name,
			Doc:  m.Doc,
		})
	}

	return &ParsedStruct{
		Fields:  fields,
		Doc:     doc,
		Methods: methods,
	}
}
