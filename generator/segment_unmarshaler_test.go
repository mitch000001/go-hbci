package generator

import (
	"bytes"
	"io/ioutil"
	"testing"

	"go/parser"
	"go/token"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestSegmentUnmarshalerGeneratorGenerate(t *testing.T) {
	testSrc := `package testsegment

import (
	"github.com/mitch000001/go-hbci/element"
)

type SegmentTest struct {
	segment.Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (s *SegmentTest) elements() []element.DataElement {
	return []element.DataElement{
		s.Abc,
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
	return nil
}
`

	generator := NewSegmentUnmarshaler("SegmentTest", "testsegment", fileSet, f)

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

	// unknown element in elements
	testSrc = `package testsegment

import (
	"github.com/mitch000001/go-hbci/element"
)

type SegmentTest struct {
	segment.Segment
	Abc *element.AlphaNumericDataElement
	Def *element.NumberDataElement
}

func (s *SegmentTest) elements() []element.DataElement {
	return []element.DataElement{
		s.Abc,
		&element.NumberDataElement{},
	}
}
`
	fileSet = token.NewFileSet()
	f, err = parser.ParseFile(fileSet, "", testSrc, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	generator = NewSegmentUnmarshaler("SegmentTest", "testsegment", fileSet, f)

	_, err = generator.Generate()

	if err != nil {
		errMessage := err.Error()
		expectedMessage := `*generator.SegmentUnmarshalerGenerator: Unsupported element in elements method: "&element.NumberDataElement{}"`
		if expectedMessage != errMessage {
			t.Logf("Expected error message to equal\n%q\n\tgot\n%q\n", expectedMessage, errMessage)
			t.Fail()
		}
	} else {
		t.Logf("Expected error, got nil\n")
		t.Fail()
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
