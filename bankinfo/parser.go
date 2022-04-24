package bankinfo

import (
	Csv "encoding/csv"
	"io"
	"strings"

	"github.com/mitch000001/go-hbci/internal"
	"github.com/pkg/errors"
	"github.com/wildducktheories/go-csv"
)

const (
	bankIdentifierHeader = "BLZ"
	bicHeader            = "BIC"
	bankInstituteHeader  = "Institut"
	versionNumberHeader  = "HBCI-Version"
	urlHeader            = "PIN/TAN-Zugang URL"
	versionNameHeader    = "Version"
	cityHeader           = "Ort"
	lastChangedHeader    = "Datum letzte Änderung"
)

const (
	bicBankIdentifier = "Bank-leitzahl"
	bicIdentifier     = "BIC"
)

// ParseBankInfos extracts all bank information from the given reader. It
// expects the reader contents to be a CSV file with ';' as separator.
func ParseBankInfos(reader io.Reader) ([]BankInfo, error) {
	CsvReader := Csv.NewReader(reader)
	CsvReader.Comma = ';'
	CsvReader.FieldsPerRecord = -1
	CsvReader.TrimLeadingSpace = true
	csvReader := csv.WithCsvReader(CsvReader, nil)
	records, err := csv.ReadAll(csvReader)
	if err != nil {
		return nil, errors.WithMessage(err, "read CSV file")
	}
	var bankInfos []BankInfo
	for _, record := range records {
		if record.Get(bankIdentifierHeader) == "" {
			internal.Debug.Printf("No BankIdentifier found for record:\n%#v\n", record.AsMap())
			continue
		}
		bankInfo := BankInfo{
			BankID:        strings.TrimSpace(record.Get(bankIdentifierHeader)),
			BIC:           strings.TrimSpace(record.Get(bicIdentifier)),
			VersionNumber: strings.TrimSpace(record.Get(versionNumberHeader)),
			URL:           strings.TrimSpace(record.Get(urlHeader)),
			VersionName:   strings.TrimSpace(record.Get(versionNameHeader)),
			Institute:     strings.TrimSpace(record.Get(bankInstituteHeader)),
			City:          strings.TrimSpace(record.Get(cityHeader)),
			LastChanged:   strings.TrimSpace(record.Get(lastChangedHeader)),
		}
		bankInfos = append(bankInfos, bankInfo)
	}
	return bankInfos, nil
}

// ParseBicData extracts all bic information from the given reader. It
// expects the reader contents to be a CSV file with ';' as separator.
func ParseBicData(reader io.Reader) ([]BicInfo, error) {
	CsvReader := Csv.NewReader(reader)
	CsvReader.Comma = ';'
	CsvReader.FieldsPerRecord = -1
	CsvReader.TrimLeadingSpace = true
	csvReader := csv.WithCsvReader(CsvReader, nil)
	records, err := csv.ReadAll(csvReader)
	if err != nil {
		return nil, err
	}
	var bicInfos []BicInfo
	for _, record := range records {
		if record.Get(bicBankIdentifier) == "" {
			internal.Debug.Printf("No BankIdentifier found for record:\n%#v\n", record.AsMap())
			continue
		}
		bic := BicInfo{
			BankID: strings.TrimSpace(record.Get(bicBankIdentifier)),
			BIC:    strings.TrimSpace(record.Get(bicIdentifier)),
		}
		bicInfos = append(bicInfos, bic)
	}
	return bicInfos, nil
}
