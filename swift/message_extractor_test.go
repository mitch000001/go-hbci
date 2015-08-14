package swift

import (
	"reflect"
	"testing"
)

func TestMessageExtractorExtract(t *testing.T) {
	test := "\r\n:20:abcde\r\n:21:def\r\n-"

	extractor := NewMessageExtractor([]byte(test))

	messages, err := extractor.Extract()

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if messages == nil {
		t.Logf("Expected messages not to be nil")
		t.Fail()
	}

	expectedMessages := [][]byte{
		[]byte("\r\n:20:abcde\r\n:21:def\r\n-"),
	}

	if !reflect.DeepEqual(expectedMessages, messages) {
		t.Logf("Expected result to equal\n%q\n\tgot\n%q\n", expectedMessages, messages)
		t.Fail()
	}
}
