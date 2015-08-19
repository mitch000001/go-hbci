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
