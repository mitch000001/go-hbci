package swift

import (
	"reflect"
	"testing"
)

func TestTagExtractorExtract(t *testing.T) {
	test := "\r\n:20:abcde\r\n:21:def\r\n-"

	extractor := NewTagExtractor([]byte(test))

	tags, err := extractor.Extract()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if tags == nil {
		t.Logf("Expected tags not to be nil")
		t.Fail()
	}

	expectedTags := [][]byte{
		[]byte(":20:abcde"),
		[]byte(":21:def"),
	}

	if !reflect.DeepEqual(expectedTags, tags) {
		t.Logf("Expected result to equal\n%q\n\tgot\n%q\n", expectedTags, tags)
		t.Fail()
	}
}
