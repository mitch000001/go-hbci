package bankinfo

import "testing"

func TestSortableBicInfosSortInterface(t *testing.T) {
	sorted := SortableBicInfos{
		BicInfo{BankId: "30000000"},
		BicInfo{BankId: "10000000"},
		BicInfo{BankId: "20000000"},
	}

	length := sorted.Len()

	if len(sorted) != length {
		t.Logf("Expected length to equal %d, got %d\n", len(sorted), length)
		t.Fail()
	}

	less := sorted.Less(0, 1)

	if less {
		t.Logf("Expected first entry not to be less than second, but was\n")
		t.Fail()
	}

	sorted.Swap(0, 1)

	if sorted[0].BankId != "10000000" {
		t.Logf("Expected first entry to have BankId '10000000', but was %q\n", sorted[0].BankId)
		t.Fail()
	}

	if sorted[1].BankId != "30000000" {
		t.Logf("Expected first entry to have BankId '30000000', but was %q\n", sorted[1].BankId)
		t.Fail()
	}
}