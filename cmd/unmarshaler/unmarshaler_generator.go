package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var segmentFlag = flag.String("segment", "", "-segment 'MyAwesomeSegment'")

func main() {
	flag.Parse()
	if *segmentFlag == "" {
		fmt.Printf("You must provide a segment to generate the unmarshaler\n")
		os.Exit(1)
	}
	filename := os.Getenv("GOFILE")
	packageName := os.Getenv("GOPACKAGE")
	fmt.Printf("Filename: %s\n", filename)
	fmt.Printf("Package: %s\n", packageName)
	absPath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Printf("Error: %T:%v\n")
	} else {
		fmt.Printf("Absolute filename: %s\n", absPath)
	}
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, filename, nil, 0)
	if err != nil {
		fmt.Println(err)
	}
	object := f.Scope.Lookup(*segmentFlag)
	visitor := &structVisitor{fileSet: fileSet}
	ast.Walk(visitor, object.Decl.(*ast.TypeSpec))

	sortedFields := sortedFields(visitor.fields)
	sort.Sort(sortedFields)

	r, _ := utf8.DecodeRuneInString(*segmentFlag)
	nameVar := string(unicode.ToLower(r))
	templObj := &segmentTemplateObject{
		Package: packageName,
		Name:    *segmentFlag,
		NameVar: nameVar,
		Fields:  sortedFields,
	}
	funcMap := map[string]interface{}{
		"plusOne": func(in int) int { return in + 1 },
	}
	t := template.Must(template.New("segment").Funcs(funcMap).Parse(segmentUnmarshalingTemplate))
	var buf bytes.Buffer
	err = t.Execute(&buf, templObj)
	if err != nil {
		fmt.Printf("Error while executing template: %T:%v\n", err, err)
		os.Exit(1)
	}
	newFileName := strings.TrimSuffix(filename, ".go") + "_unmarshaler.go"
	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Printf("Error while creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	fileSet = token.NewFileSet()
	newAstFile, err := parser.ParseFile(fileSet, newFileName, buf.String(), 0)
	if err != nil {
		fmt.Println(err)
	}
	err = printer.Fprint(file, fileSet, newAstFile)
	if err != nil {
		fmt.Println(err)
	}
	//if object != nil {
	//ast.Print(fileSet, object)
	//}
	//printer.Fprint(os.Stdout, fileSet, f)
	//ast.Print(fileSet, f)
}

type segmentTemplateObject struct {
	Package string
	Name    string
	NameVar string
	Fields  []field
	counter int
}

var segmentUnmarshalingTemplate = `package {{.Package}}

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func ({{.NameVar}} *{{.Name}}) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], {{.NameVar}})
	if err != nil {
		return err
	}
	{{.NameVar}}.Segment = seg{{ range $idx, $field := .Fields }}
	if len(elements) >= {{ plusOne $idx | plusOne }} {
		{{ $.NameVar }}.{{ $field.Name }} = &{{ $field.TypeDecl }}{}
		err = {{ $.NameVar }}.{{ $field.Name }}.UnmarshalHBCI(elements[{{ plusOne $idx }}])
		if err != nil {
			return err
		}
	}{{ end }}
	return nil
}
`

type structVisitor struct {
	fileSet *token.FileSet
	fields  []field
}

type field struct {
	Name     string
	TypeDecl string
	Line     int
}

type sortedFields []field

func (s sortedFields) Len() int           { return len(s) }
func (s sortedFields) Less(i, j int) bool { return s[i].Line < s[j].Line }
func (s sortedFields) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *structVisitor) Visit(node ast.Node) ast.Visitor {
	if structType, ok := node.(*ast.StructType); ok {
		if fields := structType.Fields; fields != nil {
			for _, f := range fields.List {
				var fieldName string
				var typeDecl string
				if names := f.Names; names != nil {
					fieldName = nodeToString(names[0], s.fileSet)
				} else {
					continue // anonymous field
				}
				if starExpr, ok := f.Type.(*ast.StarExpr); ok {
					typeDecl = nodeToString(starExpr.X, s.fileSet)
				}
				pos := f.Pos()
				position := s.fileSet.Position(pos)
				s.fields = append(s.fields, field{fieldName, typeDecl, position.Line})
			}
			return nil
		}
	}
	return s
}

func nodeToString(node ast.Node, fileSet *token.FileSet) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fileSet, node)
	return buf.String()
}
