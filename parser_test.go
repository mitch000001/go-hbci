package hbci

import (
	"reflect"
	"testing"
)

func TestParserPhase1(t *testing.T) {
	testInput := "ab??c:d+ef'"
	l := NewStringLexer("Phase1", testInput)
	p := NewParser()

	expectedAst := []*Segment{
		&Segment{
			tokens: []Token{
				elementToken{ALPHA_NUMERIC, "ab", 0},
				elementToken{ESCAPE_SEQUENCE, "??", 2},
				elementToken{ALPHA_NUMERIC, "c", 4},
				elementToken{GROUP_DATA_ELEMENT_SEPARATOR, ":", 5},
				elementToken{ALPHA_NUMERIC, "d", 6},
				elementToken{DATA_ELEMENT_SEPARATOR, "+", 7},
				elementToken{ALPHA_NUMERIC, "ef", 8},
				elementToken{SEGMENT_END_MARKER, "'", 10},
			},
			dataElements: []DataElement{
				DataElement{
					tokens: nil,
					DataElementGroup: &DataElementGroup{
						groupDataElements: []GroupDataElement{
							GroupDataElement{
								tokens: []Token{
									elementToken{ALPHA_NUMERIC, "ab", 0},
									elementToken{ESCAPE_SEQUENCE, "??", 2},
									elementToken{ALPHA_NUMERIC, "c", 4},
									elementToken{GROUP_DATA_ELEMENT_SEPARATOR, ":", 5},
								},
							},
							GroupDataElement{
								tokens: []Token{
									elementToken{ALPHA_NUMERIC, "d", 6},
									elementToken{DATA_ELEMENT_SEPARATOR, "+", 7},
								},
							},
						},
					},
				},
				DataElement{
					tokens: []Token{
						elementToken{ALPHA_NUMERIC, "ef", 8},
						elementToken{SEGMENT_END_MARKER, "'", 10},
					},
					DataElementGroup: nil,
				},
			},
		},
	}

	actualAst, err := p.Phase1(l)

	if err != nil {
		t.Logf("Expected no error, got %T: %v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedAst, actualAst) {
		t.Logf("Expected ast to equal \n%+#v \n\tgot: \n%+#v\n", expectedAst, actualAst)
		t.Fail()
	}

	type testData struct {
		input             string
		dataElements      int
		groupDataElements int
		segments          int
	}

	tests := []testData{
		{"a+b+c'", 3, 0, 1},
		{"a:b+c'", 2, 2, 1},
		{"a:b+c'd+e+f'", 5, 2, 2},
		{"a:b+c'd+e+f'", 5, 2, 2},
		{"a::b+c'd+e+f'", 5, 3, 2},
		{"a:b++c'd+e+f'", 6, 2, 2},
	}
	for _, test := range tests {
		lexer := NewStringLexer("Test Phase1", test.input)
		parser := NewParser()
		segments, err := parser.Phase1(lexer)

		if err != nil {
			t.Logf("Input: %q\n", test.input)
			t.Logf("Expected no error, got %T: %v\n", err, err)
			t.Fail()
		}

		if len(segments) != test.segments {
			t.Logf("Input: %q\n", test.input)
			t.Logf("Expected %d segments, got %d\n", test.segments, len(segments))
			t.Fail()
		}

		dataElementCount := 0
		groupDataElementCount := 0
		for _, segment := range segments {
			for _, de := range segment.dataElements {
				dataElementCount += 1
				if de.IsDataElementGroup() {
					groupDataElementCount += len(de.groupDataElements)
				}
			}
		}

		if dataElementCount != test.dataElements {
			t.Logf("Input: %q\n", test.input)
			t.Logf("Expected %d segments, got %d\n", test.dataElements, dataElementCount)
			t.Fail()
		}

		if groupDataElementCount != test.groupDataElements {
			t.Logf("Input: %q\n", test.input)
			t.Logf("Expected %d segments, got %d\n", test.groupDataElements, groupDataElementCount)
			t.Fail()
		}

	}
}
