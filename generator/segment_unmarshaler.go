package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"sort"
	"text/template"
	"unicode"
	"unicode/utf8"
)

func NewSegmentUnmarshaler(segmentName, packageName string, fileSet *token.FileSet, file *ast.File) *SegmentUnmarshalerGenerator {
	return &SegmentUnmarshalerGenerator{
		segmentName: segmentName,
		packageName: packageName,
		fileSet:     fileSet,
		file:        file,
	}
}

type SegmentUnmarshalerGenerator struct {
	segmentName string
	packageName string
	fileSet     *token.FileSet
	file        *ast.File
}

func (s *SegmentUnmarshalerGenerator) Generate() (io.Reader, error) {
	object := s.file.Scope.Lookup(s.segmentName)
	if object == nil {
		return nil, fmt.Errorf("%T: No segment with name %q found in package %q", s, s.segmentName, s.packageName)
	}
	visitor := &structVisitor{fileSet: s.fileSet}
	ast.Walk(visitor, object.Decl.(*ast.TypeSpec))

	sortedFields := sortedFields(visitor.fields)
	sort.Sort(sortedFields)

	r, _ := utf8.DecodeRuneInString(s.segmentName)
	nameVar := string(unicode.ToLower(r))
	templObj := &segmentTemplateObject{
		Package: s.packageName,
		Name:    s.segmentName,
		NameVar: nameVar,
		Fields:  sortedFields,
	}
	funcMap := map[string]interface{}{
		"plusOne": func(in int) int { return in + 1 },
	}
	t := template.Must(template.New("segment").Funcs(funcMap).Parse(segmentUnmarshalingTemplate))
	var buf bytes.Buffer
	err := t.Execute(&buf, templObj)
	if err != nil {
		return nil, fmt.Errorf("%T: Error while executing template: %v", s, err)
	}
	return &buf, nil
}

type segmentTemplateObject struct {
	Package string
	Name    string
	NameVar string
	Fields  []field
	counter int
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

type structVisitor struct {
	fileSet *token.FileSet
	fields  []field
}

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

const segmentUnmarshalingTemplate = `package {{.Package}}

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
	if len(elements) > {{ plusOne $idx }} {
		{{ $.NameVar }}.{{ $field.Name }} = &{{ $field.TypeDecl }}{}
		err = {{ $.NameVar }}.{{ $field.Name }}.UnmarshalHBCI(elements[{{ plusOne $idx }}])
		if err != nil {
			return err
		}
	}{{ end }}
	return nil
}
`
