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

	generator := NewSegmentUnmarshaler(SegmentIdentifier{Name: "TestSegment", InterfaceName: "Segment"}, "test_files", fileSet, f)

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

	generator = NewSegmentUnmarshaler(SegmentIdentifier{Name: "TestSegmentUnknownElement"}, "test_files", fileSet, f)

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
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "test_files/versioned_test_segment.go", nil, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	expectedSrc, err := ioutil.ReadFile("test_files/versioned_test_segment_unmarshaler.go")

	segment := SegmentIdentifier{
		Name:          "VersionedTestSegment",
		InterfaceName: "BankSegment",
		Versions: []SegmentIdentifier{
			{
				Name:          "VersionedTestSegmentV1",
				Version:       1,
				InterfaceName: "Segment",
			},
		},
	}

	generator := NewVersionedSegmentUnmarshaler(segment, "test_files", fileSet, f)

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

	// custom interface
	fileSet = token.NewFileSet()
	f, err = parser.ParseFile(fileSet, "test_files/versioned_test_segment_custom_interface.go", nil, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	expectedSrc, err = ioutil.ReadFile("test_files/versioned_test_segment_custom_interface_unmarshaler.go")

	segment = SegmentIdentifier{
		Name:          "VersionedTestSegmentCustomInterface",
		InterfaceName: "versionedTestSegmentCustomInterface",
		Versions: []SegmentIdentifier{
			{
				Name:          "VersionedTestSegmentCustomInterfaceV1",
				Version:       1,
				InterfaceName: "Segment",
			},
		},
	}

	generator = NewVersionedSegmentUnmarshaler(segment, "test_files", fileSet, f)

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
		if !bytes.Equal(expectedSrc, generatedSourcebytes) {
			diffs := diffmatchpatch.New().DiffMain(string(expectedSrc), string(generatedSourcebytes), true)
			t.Logf("Expected generated sources to equal\n%s\n\tgot\n%s\n", expectedSrc, generatedSourcebytes)
			t.Logf("Diff: \n%s\n", diffPrettyPrint(diffs))
			t.Fail()
		}
	}

	// multiple versions
	fileSet = token.NewFileSet()
	f, err = parser.ParseFile(fileSet, "test_files/multiple_versioned_test_segment.go", nil, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	expectedSrc, err = ioutil.ReadFile("test_files/multiple_versioned_test_segment_unmarshaler.go")

	segment = SegmentIdentifier{
		Name:          "MultipleVersionedTestSegment",
		InterfaceName: "BankSegment",
		Versions: []SegmentIdentifier{
			{
				Name:          "MultipleVersionedTestSegmentV1",
				Version:       1,
				InterfaceName: "Segment",
			},
			{
				Name:          "MultipleVersionedTestSegmentV2",
				Version:       2,
				InterfaceName: "Segment",
			},
		},
	}

	generator = NewVersionedSegmentUnmarshaler(segment, "test_files", fileSet, f)

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
		if !bytes.Equal(expectedSrc, generatedSourcebytes) {
			diffs := diffmatchpatch.New().DiffMain(string(expectedSrc), string(generatedSourcebytes), true)
			t.Logf("Expected generated sources to equal\n%s\n\tgot\n%s\n", expectedSrc, generatedSourcebytes)
			t.Logf("Diff: \n%s\n", diffPrettyPrint(diffs))
			t.Fail()
		}
	}

	// multiple versions custom interfaces
	fileSet = token.NewFileSet()
	f, err = parser.ParseFile(fileSet, "test_files/multiple_versioned_test_segment_custom_interfaces.go", nil, 0)
	if err != nil {
		t.Logf("Error while parsing source: %T:%v\n", err, err)
		t.FailNow()
	}

	expectedSrc, err = ioutil.ReadFile("test_files/multiple_versioned_test_segment_custom_interfaces_unmarshaler.go")

	segment = SegmentIdentifier{
		Name:          "MultipleVersionedTestSegmentCustomInterfaces",
		InterfaceName: "BankSegment",
		Versions: []SegmentIdentifier{
			{
				Name:          "MultipleVersionedTestSegmentCustomInterfacesV1",
				Version:       1,
				InterfaceName: "versionInterface1",
			},
			{
				Name:          "MultipleVersionedTestSegmentCustomInterfacesV2",
				Version:       2,
				InterfaceName: "versionInterface2",
			},
		},
	}

	generator = NewVersionedSegmentUnmarshaler(segment, "test_files", fileSet, f)

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
		if !bytes.Equal(expectedSrc, generatedSourcebytes) {
			diffs := diffmatchpatch.New().DiffMain(string(expectedSrc), string(generatedSourcebytes), true)
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
