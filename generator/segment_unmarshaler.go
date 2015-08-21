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
	fieldExtractor := &fieldExtractor{
		segment: s.segment,
		file:    s.file,
		fileSet: s.fileSet,
	}
	sortedFields, err := fieldExtractor.extractFields()
	if err != nil {
		return nil, err
	}

	r, _ := utf8.DecodeRuneInString(s.segment.Name)
	nameVar := string(unicode.ToLower(r))
	templObj := &segmentTemplateObject{
		Package:       s.packageName,
		Name:          s.segment.Name,
		NameVar:       nameVar,
		InterfaceName: s.segment.InterfaceName,
		Fields:        sortedFields,
	}
	executor := &segmentTemplateExecutor{templObj}
	return executor.execute()
}

func NewVersionedSegmentUnmarshaler(segment SegmentIdentifier, packageName string, fileSet *token.FileSet, file *ast.File) *VersionedSegmentUnmarshalerGenerator {
	return &VersionedSegmentUnmarshalerGenerator{
		segment:     segment,
		packageName: packageName,
		fileSet:     fileSet,
		file:        file,
	}
}

type VersionedSegmentUnmarshalerGenerator struct {
	segment     SegmentIdentifier
	packageName string
	fileSet     *token.FileSet
	file        *ast.File
}

func (v *VersionedSegmentUnmarshalerGenerator) Generate() (io.Reader, error) {
	var versionedTemplateObjects []*segmentTemplateObject
	for _, version := range v.segment.Versions {
		fieldExtractor := &fieldExtractor{
			segment: version,
			file:    v.file,
			fileSet: v.fileSet,
		}
		sortedFields, err := fieldExtractor.extractFields()
		if err != nil {
			return nil, err
		}

		r, _ := utf8.DecodeRuneInString(version.Name)
		nameVar := string(unicode.ToLower(r))
		templObj := &segmentTemplateObject{
			Package:       v.packageName,
			Name:          version.Name,
			NameVar:       nameVar,
			InterfaceName: version.InterfaceName,
			Version:       version.Version,
			Fields:        sortedFields,
		}
		versionedTemplateObjects = append(versionedTemplateObjects, templObj)
	}
	r, _ := utf8.DecodeRuneInString(v.segment.Name)
	nameVar := string(unicode.ToLower(r))
	templObj := &segmentTemplateObject{
		Package:       v.packageName,
		Name:          v.segment.Name,
		NameVar:       nameVar,
		InterfaceName: v.segment.InterfaceName,
	}
	segmentTemplObj := &versionedSegmentTemplateObject{
		segmentTemplateObject: templObj,
		SegmentVersions:       versionedTemplateObjects,
	}

	executor := &versionedSegmentTemplateExecutor{segmentTemplObj}
	return executor.execute()
}

type segmentTemplateExecutor struct {
	templateObject *segmentTemplateObject
}

func (s *segmentTemplateExecutor) execute() (io.Reader, error) {
	funcMap := map[string]interface{}{
		"plusOne": func(in int) int { return in + 1 },
	}
	t, err := template.New("executor").Funcs(funcMap).Parse(segmentExecutorTemplate)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template: %v", err)
	}
	t, err = t.Parse(segmentUnmarshalingTemplate)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template: %v", err)
	}
	t, err = t.Parse(packageDeclTemplate)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template: %v", err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, s.templateObject)
	if err != nil {
		return nil, fmt.Errorf("%T: Error while executing template: %v", s, err)
	}
	return &buf, nil
}

type segmentTemplateObject struct {
	Package       string
	Name          string
	NameVar       string
	InterfaceName string
	Version       int
	Fields        []field
	counter       int
}

type fieldExtractor struct {
	file    *ast.File
	fileSet *token.FileSet
	segment SegmentIdentifier
}

func (f *fieldExtractor) extractFields() ([]field, error) {
	object := f.file.Scope.Lookup(f.segment.Name)
	if object == nil {
		return nil, fmt.Errorf("No segment with name %q found", f.segment.Name)
	}
	elemVisitor := &elementVisitor{fileSet: f.fileSet, object: object, receiverName: f.segment.Name}
	ast.Walk(elemVisitor, f.file)
	if elemVisitor.err != nil {
		return nil, elemVisitor.err
	}

	sortedFields := sortedFields(elemVisitor.fields)
	sort.Sort(sortedFields)
	return sortedFields, nil
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

const segmentExecutorTemplate = `{{template "package_declaration" .}}
{{template "segment" .}}
`

const packageDeclTemplate = `{{define "package_declaration"}}package {{.Package}}

import (
	"bytes"
	"fmt"

	"github.com/mitch000001/go-hbci/element"
){{end}}
`

const segmentUnmarshalingTemplate = `{{define "segment"}}
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
	{{.NameVar}}.{{.InterfaceName}} = seg{{ range $idx, $field := .Fields }}
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
}{{end}}
`

type versionedSegmentTemplateObject struct {
	*segmentTemplateObject
	SegmentVersions []*segmentTemplateObject
}

type versionedSegmentTemplateExecutor struct {
	templateObject *versionedSegmentTemplateObject
}

func (v *versionedSegmentTemplateExecutor) execute() (io.Reader, error) {
	funcMap := map[string]interface{}{
		"plusOne": func(in int) int { return in + 1 },
	}
	t, err := template.New("executor").Funcs(funcMap).Parse(versionedSegmentExecutorTemplate)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template: %v", err)
	}
	t, err = t.Parse(versionedSegmentUnmarshalingTemplate)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template: %v", err)
	}
	t, err = t.Parse(segmentUnmarshalingTemplate)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template: %v", err)
	}
	t, err = t.Parse(packageDeclTemplate)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing template: %v", err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, v.templateObject)
	if err != nil {
		return nil, fmt.Errorf("Error while executing template: %v", err)
	}
	return &buf, nil
}

const versionedSegmentExecutorTemplate = `{{template "package_declaration" .}}
{{template "versioned_segment" .}}
`

const versionedSegmentUnmarshalingTemplate = `
{{define "versioned_segment"}}
func ({{.NameVar}} *{{.Name}}) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment {{.InterfaceName}}
	switch header.Version.Val() {
	{{range $version := .SegmentVersions}}case {{$version.Version}}:{{with $versionName := $version.Name}}
		segment = &{{$versionName}}{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}{{end}}
	{{end}}default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	{{.NameVar}}.{{.InterfaceName}} = segment
	return nil
}{{ range $versioned := .SegmentVersions }}
{{template "segment" $versioned}}{{end}}{{end}}
`
