package generator

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

	expectedSrc, err := os.ReadFile("test_files/test_segment_unmarshaler.go")
	if err != nil {
		t.Logf("Error reading test fixtures: %v", err)
		t.FailNow()
	}

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
		generatedSourcebytes, err := io.ReadAll(reader)
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
	testGenerator := func(t *testing.T, fixtureFile string, expectedFixtureFile string, segmentIdentifier SegmentIdentifier) {
		fileSet := token.NewFileSet()
		f, err := parser.ParseFile(fileSet, fixtureFile, nil, 0)
		if err != nil {
			t.Logf("Error while parsing source: %T:%v\n", err, err)
			t.FailNow()
		}

		expectedSrc, err := os.ReadFile(expectedFixtureFile)
		if err != nil {
			t.Logf("Error reading test fixtures: %v", err)
			t.FailNow()
		}

		generator := NewVersionedSegmentUnmarshaler(segmentIdentifier, "test_files", fileSet, f)

		reader, err := generator.Generate()

		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		if reader == nil {
			t.Logf("Expected reader not to be nil")
			t.Fail()
		} else {
			generatedSourcebytes, err := io.ReadAll(reader)
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

	t.Run("versioned segment", func(t *testing.T) {
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
		testGenerator(t, "test_files/versioned_test_segment.go", "test_files/versioned_test_segment_unmarshaler.go", segment)
	})

	t.Run("custom interface", func(t *testing.T) {
		segment := SegmentIdentifier{
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
		testGenerator(t, "test_files/versioned_test_segment_custom_interface.go", "test_files/versioned_test_segment_custom_interface_unmarshaler.go", segment)
	})

	t.Run("multiple versions", func(t *testing.T) {
		segment := SegmentIdentifier{
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
		testGenerator(t, "test_files/multiple_versioned_test_segment.go", "test_files/multiple_versioned_test_segment_unmarshaler.go", segment)
	})

	t.Run("multiple versions custom interfaces", func(t *testing.T) {
		segment := SegmentIdentifier{
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
		testGenerator(t, "test_files/multiple_versioned_test_segment_custom_interfaces.go", "test_files/multiple_versioned_test_segment_custom_interfaces_unmarshaler.go", segment)
	})
}

func diffPrettyPrint(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString("\u2588>>> + >>> ")
			fmt.Fprintf(&buff, "%q", text)
			// buff.WriteString(text)
			buff.WriteString(" <<< + <<<\u2588")
		case diffmatchpatch.DiffDelete:
			buff.WriteString("\u2588>>> - >>> ")
			fmt.Fprintf(&buff, "%q", text)
			// buff.WriteString(text)
			buff.WriteString(" <<< - <<<\u2588")
		case diffmatchpatch.DiffEqual:
			buff.WriteString(text)
		}
	}
	return buff.String()
}
