package bankinfo

import (
	Csv "encoding/csv"
	"io"

	"github.com/wildducktheories/go-csv"
)

const (
	BANK_IDENTIFIER       = "BLZ"
	VERSION_NUMBER_HEADER = "HBCI-Version"
	URL_HEADER            = "PIN/TAN-Zugang URL"
	VERSION_HEADER        = "Version"
)

type Parser struct {
}

func (p Parser) Parse(reader io.Reader) ([]BankInfo, error) {
	CsvReader := Csv.NewReader(reader)
	CsvReader.Comma = ';'
	CsvReader.FieldsPerRecord = -1
	CsvReader.TrimLeadingSpace = true
	csvReader, err := csv.WithCsvReader(CsvReader)
	if err != nil {
		return nil, err
	}
	records, err := csv.ReadAll(csvReader)
	var bankInfos []BankInfo
	for _, record := range records {
		if record.Get(BANK_IDENTIFIER) == "" {
			continue
		}
		bankInfo := BankInfo{
			BankId:        record.Get(BANK_IDENTIFIER),
			VersionNumber: record.Get(VERSION_NUMBER_HEADER),
			URL:           record.Get(URL_HEADER),
			VersionString: record.Get(VERSION_HEADER),
		}
		bankInfos = append(bankInfos, bankInfo)
	}
	return bankInfos, nil
}
