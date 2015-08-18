package swift

import (
	"reflect"
	"testing"
)

func TestTagExtractorExtract(t *testing.T) {
	tests := []struct {
		marshaled string
		result    [][]byte
	}{
		{
			"\r\n:20:abcde\r\n:21:def\r\n-",
			[][]byte{
				[]byte(":20:abcde"),
				[]byte(":21:def"),
			},
		},
		{
			"\r\n:20:abc\r\nde\r\n:21:def\r\n-",
			[][]byte{
				[]byte(":20:abc\r\nde"),
				[]byte(":21:def"),
			},
		},
	}
	for _, test := range tests {
		extractor := NewTagExtractor([]byte(test.marshaled))

		tags, err := extractor.Extract()

		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		if tags == nil {
			t.Logf("Expected tags not to be nil")
			t.Fail()
		}

		if !reflect.DeepEqual(test.result, tags) {
			t.Logf("Expected result to equal\n%q\n\tgot\n%q\n", test.result, tags)
			t.Fail()
		}
	}
}
