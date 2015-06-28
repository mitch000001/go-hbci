package hbci

import "testing"

func TestSegmentHeaderGroupDataElements(t *testing.T) {
	header := NewSegmentHeader("abc", 1, 2)

	elements := header.GroupDataElements()

	expectedLength := 4

	if expectedLength != len(elements) {
		t.Logf("Expected %d GroupDataElements, got %d\n", expectedLength, len(elements))
		t.Fail()
	}

	header = NewReferencingSegmentHeader("abc", 1, 2, 3)

	elements = header.GroupDataElements()

	expectedLength = 4

	if expectedLength != len(elements) {
		t.Logf("Expected %d GroupDataElements, got %d\n", expectedLength, len(elements))
		t.Fail()
	}
}

func TestSegmentHeaderString(t *testing.T) {
	header := NewSegmentHeader("abc", 1, 2)

	expectedString := "abc:1:2:"

	actualString := header.String()

	if expectedString != actualString {
		t.Logf("Expected String() to equal %q, was %q\n", expectedString, actualString)
		t.Fail()
	}
}

func TestSegmentHeaderType(t *testing.T) {
	header := NewSegmentHeader("abc", 1, 2)

	expectedType := SegmentHeaderDEG

	actualType := header.Type()

	if expectedType != actualType {
		t.Logf("Expected Type() to return %v, got %v\n", expectedType, actualType)
		t.Fail()
	}
}
