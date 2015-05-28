package hbci

import (
	"reflect"
	"testing"
)

func TestParserPhase1(t *testing.T) {
	testInput := "ab??c:d+ef'"
	l := NewLexer("Phase1", testInput)
	p := NewParser()

	expectedAst := []*Segment{
		&Segment{
			tokens: []Token{
				Token{ALPHA_NUMERIC, "ab", 0},
				Token{ESCAPE_SEQUENCE, "??", 2},
				Token{ALPHA_NUMERIC, "c", 4},
				Token{GROUP_DATA_ELEMENT_SEPARATOR, ":", 5},
				Token{ALPHA_NUMERIC, "d", 6},
				Token{DATA_ELEMENT_SEPARATOR, "+", 7},
				Token{ALPHA_NUMERIC, "ef", 8},
				Token{SEGMENT_END_MARKER, "'", 10},
			},
			dataElements: []DataElement{
				DataElement{
					tokens: nil,
					DataElementGroup: &DataElementGroup{
						groupDataElements: []GroupDataElement{
							GroupDataElement{
								tokens: []Token{
									Token{ALPHA_NUMERIC, "ab", 0},
									Token{ESCAPE_SEQUENCE, "??", 2},
									Token{ALPHA_NUMERIC, "c", 4},
									Token{GROUP_DATA_ELEMENT_SEPARATOR, ":", 5},
								},
							},
							GroupDataElement{
								tokens: []Token{
									Token{ALPHA_NUMERIC, "d", 6},
									Token{DATA_ELEMENT_SEPARATOR, "+", 7},
								},
							},
						},
					},
				},
				DataElement{
					tokens: []Token{
						Token{ALPHA_NUMERIC, "ef", 8},
						Token{SEGMENT_END_MARKER, "'", 10},
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
}
