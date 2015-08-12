package element

import "testing"

func TestReferencingMessageDataElementUnmarshalHBCI(t *testing.T) {
	test := "abcde:1"

	expected := NewReferencingMessage("abcde", 1).String()

	element := &ReferencingMessageDataElement{}

	err := element.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	actual := element.String()

	if expected != actual {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}
