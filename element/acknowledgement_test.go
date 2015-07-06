package element

import (
	"reflect"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
)

func TestAcknowledgementDataElementUnmarshalHBCI(t *testing.T) {
	test := "0300:7,2:Syntaxerror:test1"

	acknowledgement := &AcknowledgementDataElement{}

	err := acknowledgement.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := NewAcknowledgement(domain.Acknowledgement{300, "7,2", "Syntaxerror", []string{"test1"}})

	if !reflect.DeepEqual(expected, acknowledgement) {
		t.Logf("Expected\n%q\n\tgot\n%q\n", expected, acknowledgement)
		t.Fail()
	}
}

func TestParamsDataElementUnmarshalHBCI(t *testing.T) {
	test := "test1:test2:test3"

	params := &ParamsDataElement{}

	err := params.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := NewParams(10, 10, "test1", "test2", "test3")

	if !reflect.DeepEqual(expected, params) {
		t.Logf("Expected\n%q\n\tgot\n%q\n", expected, params)
		t.Fail()
	}
}
