package element

import (
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

	expected := NewAcknowledgement(domain.NewMessageAcknowledgement(300, "7,2", "Syntaxerror", []string{"test1"})).String()
	actual := acknowledgement.String()

	if expected != actual {
		t.Logf("Expected\n%q\n\tgot\n%q\n", expected, actual)
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

	expected := NewParams(10, 10, "test1", "test2", "test3").String()
	actual := params.String()

	if expected != actual {
		t.Logf("Expected\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}
