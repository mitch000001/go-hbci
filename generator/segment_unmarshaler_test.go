package generator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"go/parser"
	"go/token"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestSegmentUnmarshalerGeneratorGenerate(t *testing.T) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "test_files/test_segment.go", nil, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	expectedSrc, err := ioutil.ReadFile("test_files/test_segment_unmarshaler.go")

	generator := NewSegmentUnmarshaler(SegmentIdentifier{Name: "TestSegment"}, "test_files", fileSet, f)

	reader, err := generator.Generate()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if reader == nil {
		t.Logf("Expected reader not to be nil")
		t.Fail()
	} else {
		generatedSourcebytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Logf("Error while parsing source: %T:%v\n", err, err)
			t.FailNow()
		}
		if !bytes.Equal(expectedSrc, generatedSourcebytes) {
			diffs := diffmatchpatch.New().DiffMain(string(expectedSrc), string(generatedSourcebytes), true)
			t.Logf("Expected generated sources to equal\n%s\n\tgot\n%s\n", expectedSrc, generatedSourcebytes)
			t.Logf("Diff: \n%s\n", diffPrettyPrint(diffs))
			t.Fail()
		}
	}

	// unknown element in elements
	fileSet = token.NewFileSet()
	f, err = parser.ParseFile(fileSet, "test_files/test_segment_unknown_element.go", nil, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	generator = NewSegmentUnmarshaler(SegmentIdentifier{Name: "TestSegment"}, "test_files", fileSet, f)

	_, err = generator.Generate()

	if err != nil {
		errMessage := err.Error()
		expectedMessage := `Unsupported element in elements method: "&element.NumberDataElement{}"`
		if expectedMessage != errMessage {
			t.Logf("Expected error message to equal\n%q\n\tgot\n%q\n", expectedMessage, errMessage)
			t.Fail()
		}
	} else {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	}
}

func TestVersionedSegmentUnmarshalerGeneratorGenerate(t *testing.T) {
	testSrc := `package testsegment

import (
	"github.com/mitch000001/go-hbci/element"
)

type SegmentTest struct {
	segment.Segment
}

type SegmentTestV1 struct {
	segment.Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (s *SegmentTestV1) elements() []element.DataElement {
	return []element.DataElement{
		s.Abc,
		s.Def,
	}
}
`
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "", testSrc, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	expectedSrc := `package testsegment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (s *SegmentTest) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment Segment
	switch header.Version.Val() {
	case 1:
		segment = &SegmentTestV1{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	s.Segment = segment
	return nil
}

func (s *SegmentTestV1) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], s)
	if err != nil {
		return err
	}
	s.Segment = seg
	if len(elements) > 1 {
		s.Abc = &element.AlphaNumericDataElement{}
		err = s.Abc.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 {
		s.Def = &element.NumberDataElement{}
		err = s.Def.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	return nil
}
`

	segment := SegmentIdentifier{
		Name:          "SegmentTest",
		InterfaceName: "Segment",
		Versions: []SegmentIdentifier{
			{
				Name:    "SegmentTestV1",
				Version: 1,
			},
		},
	}

	generator := NewVersionedSegmentUnmarshaler(segment, "testsegment", fileSet, f)

	reader, err := generator.Generate()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if reader == nil {
		t.Logf("Expected reader not to be nil")
		t.Fail()
	} else {
		generatedSourcebytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Logf("Error while parsing source: %T:%v\n", err, err)
			t.FailNow()
		}
		if !bytes.Equal([]byte(expectedSrc), generatedSourcebytes) {
			diffs := diffmatchpatch.New().DiffMain(expectedSrc, string(generatedSourcebytes), true)
			t.Logf("Expected generated sources to equal\n%s\n\tgot\n%s\n", expectedSrc, generatedSourcebytes)
			t.Logf("Diff: \n%s\n", diffPrettyPrint(diffs))
			t.Fail()
		}
	}

	// multiple versions
	testSrc = `package testsegment

import (
	"github.com/mitch000001/go-hbci/element"
)

type SegmentTest struct {
	segment.Segment
}

type SegmentTestV1 struct {
	segment.Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (s *SegmentTestV1) elements() []element.DataElement {
	return []element.DataElement{
		s.Abc,
		s.Def,
	}
}

type SegmentTestV2 struct {
	segment.Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (s *SegmentTestV2) elements() []element.DataElement {
	return []element.DataElement{
		s.Abc,
		s.Def,
	}
}
`
	fileSet = token.NewFileSet()
	f, err = parser.ParseFile(fileSet, "", testSrc, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	expectedSrc = `package testsegment

import (
	"fmt"

	"github.com/mitch000001/go-hbci/element"
)

func (s *SegmentTest) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	header := &element.SegmentHeader{}
	err = header.UnmarshalHBCI(elements[0])
	if err != nil {
		return err
	}
	var segment Segment
	switch header.Version.Val() {
	case 1:
		segment = &SegmentTestV1{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	case 2:
		segment = &SegmentTestV2{}
		err = segment.UnmarshalHBCI(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown segment version: %d", header.Version.Val())
	}
	s.Segment = segment
	return nil
}

func (s *SegmentTestV1) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], s)
	if err != nil {
		return err
	}
	s.Segment = seg
	if len(elements) > 1 {
		s.Abc = &element.AlphaNumericDataElement{}
		err = s.Abc.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 {
		s.Def = &element.NumberDataElement{}
		err = s.Def.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SegmentTestV2) UnmarshalHBCI(value []byte) error {
	elements, err := ExtractElements(value)
	if err != nil {
		return err
	}
	if len(elements) == 0 {
		return fmt.Errorf("Malformed marshaled value")
	}
	seg, err := SegmentFromHeaderBytes(elements[0], s)
	if err != nil {
		return err
	}
	s.Segment = seg
	if len(elements) > 1 {
		s.Abc = &element.AlphaNumericDataElement{}
		err = s.Abc.UnmarshalHBCI(elements[1])
		if err != nil {
			return err
		}
	}
	if len(elements) > 2 {
		s.Def = &element.NumberDataElement{}
		err = s.Def.UnmarshalHBCI(elements[2])
		if err != nil {
			return err
		}
	}
	return nil
}
`

	segment = SegmentIdentifier{
		Name:          "SegmentTest",
		InterfaceName: "Segment",
		Versions: []SegmentIdentifier{
			{
				Name:    "SegmentTestV1",
				Version: 1,
			},
			{
				Name:    "SegmentTestV2",
				Version: 2,
			},
		},
	}

	generator = NewVersionedSegmentUnmarshaler(segment, "testsegment", fileSet, f)

	reader, err = generator.Generate()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if reader == nil {
		t.Logf("Expected reader not to be nil")
		t.Fail()
	} else {
		generatedSourcebytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Logf("Error while parsing source: %T:%v\n", err, err)
			t.FailNow()
		}
		if !bytes.Equal([]byte(expectedSrc), generatedSourcebytes) {
			diffs := diffmatchpatch.New().DiffMain(expectedSrc, string(generatedSourcebytes), true)
			t.Logf("Expected generated sources to equal\n%s\n\tgot\n%s\n", expectedSrc, generatedSourcebytes)
			t.Logf("Diff: \n%s\n", diffPrettyPrint(diffs))
			t.Fail()
		}
	}

}

func diffPrettyPrint(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString("\u2588>>> + >>> ")
			buff.WriteString(text)
			buff.WriteString(" <<< + <<<\u2588")
		case diffmatchpatch.DiffDelete:
			buff.WriteString("\u2588>>> - >>> ")
			buff.WriteString(text)
			buff.WriteString(" <<< - <<<\u2588")
		case diffmatchpatch.DiffEqual:
			buff.WriteString(text)
		}
	}
	return buff.String()
}
