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

func NewSegmentUnmarshaler(segment SegmentIdentifier, packageName string, fileSet *token.FileSet, file *ast.File) *SegmentUnmarshalerGenerator {
	return &SegmentUnmarshalerGenerator{
		segment:     segment,
		packageName: packageName,
		fileSet:     fileSet,
		file:        file,
	}
}

type SegmentUnmarshalerGenerator struct {
	segment     SegmentIdentifier
	packageName string
	fileSet     *token.FileSet
	file        *ast.File
}

func (s *SegmentUnmarshalerGenerator) Generate() (io.Reader, error) {
	sortedFields, err := s.extractFields()
	if err != nil {
		return nil, err
	}

	r, _ := utf8.DecodeRuneInString(s.segment.Name)
	nameVar := string(unicode.ToLower(r))
	templObj := &segmentTemplateObject{
		Package: s.packageName,
		Name:    s.segment.Name,
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
		return nil, fmt.Errorf("%T: Error while executing template: %v", s, err)
	}
	return &buf, nil
}

func (s *SegmentUnmarshalerGenerator) extractFields() ([]field, error) {
	object := s.file.Scope.Lookup(s.segment.Name)
	if object == nil {
		return nil, fmt.Errorf("%T: No segment with name %q found in package %q", s, s.segment.Name, s.packageName)
	}
	elemVisitor := &elementVisitor{fileSet: s.fileSet, object: object, receiverName: s.segment.Name}
	ast.Walk(elemVisitor, s.file)
	if elemVisitor.err != nil {
		return nil, fmt.Errorf("%T: %v", s, elemVisitor.err)
	}

	sortedFields := sortedFields(elemVisitor.fields)
	sort.Sort(sortedFields)
	return sortedFields, nil
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

const segmentElementsMethodIdentifier = "elements"
const segmentElementsMethodReturnType = "[]element.DataElement"

type elementVisitor struct {
	receiverName string
	fileSet      *token.FileSet
	object       *ast.Object
	fields       []field
	err          error
}

func (e *elementVisitor) Visit(node ast.Node) ast.Visitor {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		if e.isElementsMethod(funcDecl) {
			bodyStatements := funcDecl.Body.List
			if len(bodyStatements) == 1 {
				if ret, ok := bodyStatements[0].(*ast.ReturnStmt); ok {
					if len(ret.Results) == 1 {
						if res, ok := ret.Results[0].(*ast.CompositeLit); ok {
							resType := nodeToString(res.Type, e.fileSet)
							if resType == segmentElementsMethodReturnType {
								for _, element := range res.Elts {
									if sel, ok := element.(*ast.SelectorExpr); ok {
										pos := sel.Pos()
										position := e.fileSet.Position(pos)
										elemField := field{Name: sel.Sel.Name, Line: position.Line}
										if _, ok := sel.X.(*ast.Ident); ok {
											fieldVisitor := &structFieldVisitor{fileSet: e.fileSet, fieldName: sel.Sel.Name}
											ast.Walk(fieldVisitor, e.object.Decl.(*ast.TypeSpec))
											if fieldVisitor.err != nil {
												e.err = fieldVisitor.err
												return nil
											}
											elemField.TypeDecl = fieldVisitor.fieldTypeDecl
										} else {
											e.err = fmt.Errorf("Unexpected Selector: %q", nodeToString(sel, e.fileSet))
											return nil
										}
										if elemField.TypeDecl == "" {
											e.err = fmt.Errorf("No type declaration found for element %q", sel.Sel.Name)
											return nil
										}
										e.fields = append(e.fields, elemField)
									} else {
										e.err = fmt.Errorf("Unsupported element in elements method: %q", nodeToString(element, e.fileSet))
										return nil
									}
								}
								return nil
							}
						}
					}
				}
			}
		}
	}
	return e
}

func (e *elementVisitor) isElementsMethod(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Name.Name == segmentElementsMethodIdentifier {
		if funcDecl.Recv != nil {
			if fields := funcDecl.Recv.List; fields != nil {
				if len(fields) == 1 {
					if typ, ok := fields[0].Type.(*ast.StarExpr); ok {
						if ident, ok := typ.X.(*ast.Ident); ok {
							if ident.Name == e.receiverName {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

type structFieldVisitor struct {
	fileSet       *token.FileSet
	fieldName     string
	fieldTypeDecl string
	err           error
}

func (s *structFieldVisitor) Visit(node ast.Node) ast.Visitor {
	if structType, ok := node.(*ast.StructType); ok {
		if fields := structType.Fields; fields != nil {
			for _, f := range fields.List {
				var fieldName string
				if names := f.Names; names != nil {
					fieldName = nodeToString(names[0], s.fileSet)
					if fieldName != s.fieldName {
						continue // not the field we're looking for
					}
				} else {
					continue // anonymous field
				}
				if starExpr, ok := f.Type.(*ast.StarExpr); ok {
					s.fieldTypeDecl = nodeToString(starExpr.X, s.fileSet)
				} else {
					s.err = fmt.Errorf("Unexpected type found: %q", nodeToString(f.Type, s.fileSet))
					return nil
				}
			}
			return nil
		}
	}
	return s
}

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
	"bytes"
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
	if len(elements) > {{ plusOne $idx }} && len(elements[{{ plusOne $idx}}]) > 0 {
		{{ $.NameVar }}.{{ $field.Name }} = &{{ $field.TypeDecl }}{}
		{{if len $.Fields | eq (plusOne $idx)}}if len(elements)+1 > {{plusOne $idx}} {
			err = {{ $.NameVar }}.{{ $field.Name }}.UnmarshalHBCI(bytes.Join(elements[{{ plusOne $idx }}:], []byte("+")))
		} else {
			err = {{ $.NameVar }}.{{ $field.Name }}.UnmarshalHBCI(elements[{{ plusOne $idx }}])
		}{{else}}err = {{ $.NameVar }}.{{ $field.Name }}.UnmarshalHBCI(elements[{{ plusOne $idx }}]){{end}}
		if err != nil {
			return err
		}
	}{{ end }}
	return nil
}
`
