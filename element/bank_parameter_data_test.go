package element

import "testing"

func TestPinTanBusinessTransactionParametersUnmarshalHBCI(t *testing.T) {
	test := "HKSAL:N:HKUEB:J"

	param1 := &PinTanBusinessTransactionParameter{
		SegmentID: NewAlphaNumeric("HKSAL", 6),
		NeedsTAN:  NewBoolean(false),
	}
	param1.DataElement = NewGroupDataElementGroup(pinTanBusinessTransactionParameterGDEG, 2, param1)
	param2 := &PinTanBusinessTransactionParameter{
		SegmentID: NewAlphaNumeric("HKUEB", 6),
		NeedsTAN:  NewBoolean(true),
	}
	param2.DataElement = NewGroupDataElementGroup(pinTanBusinessTransactionParameterGDEG, 2, param2)
	expectedElement := &PinTanBusinessTransactionParameters{}
	expectedElement.arrayElementGroup = newArrayElementGroup(pinTanBusinessTransactionParameterGDEG, 2, 2, []DataElement{param1, param2})

	expected := expectedElement.String()

	element := &PinTanBusinessTransactionParameters{}

	err := element.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	actual := element.String()

	if actual != expected {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}

func TestPinTanBusinessTransactionParameterUnmarshalHBCI(t *testing.T) {
	test := "HKSAL:N"

	expectedElement := &PinTanBusinessTransactionParameter{
		SegmentID: NewAlphaNumeric("HKSAL", 6),
		NeedsTAN:  NewBoolean(false),
	}
	expectedElement.DataElement = NewGroupDataElementGroup(pinTanBusinessTransactionParameterGDEG, 2, expectedElement)

	expected := expectedElement.String()

	element := &PinTanBusinessTransactionParameter{}

	err := element.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	actual := element.String()

	if actual != expected {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}
