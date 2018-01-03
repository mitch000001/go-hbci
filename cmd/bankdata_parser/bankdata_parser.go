package main

import (
	"bytes"
	"flag"
	"go/format"
	"io"
	"log"
	"os"
	"sort"
	"text/template"

	"golang.org/x/text/encoding/charmap"

	"github.com/mitch000001/go-hbci/bankinfo"
	"github.com/pkg/errors"
)

func main() {
	flag.Parse()
	bankdataFiles := flag.Args()
	if len(bankdataFiles) == 0 {
		log.Fatal("No file provided. Exiting...")
		os.Exit(1)
	}

	var bankInfos []bankinfo.BankInfo
	decoder := charmap.ISO8859_1.NewDecoder()
	for _, bankdata := range bankdataFiles {
		file, err := os.Open(bankdata)
		if err != nil {
			log.Fatalf("Cannot open file: %q", bankdata)
			os.Exit(1)
		}
		infos, err := bankinfo.ParseBankInfos(decoder.Reader(file))
		if err != nil {
			log.Fatalf("Parse error: %q", err)
			os.Exit(1)
		}
		bankInfos = append(bankInfos, infos...)
	}
	sort.Sort(bankinfo.SortableBankInfos(bankInfos))
	data, err := writeDataToGoFile(bankInfos)
	if err != nil {
		log.Fatalf("Error while writing generated source: %v", err)
		os.Exit(1)
	}
	goFile, err := os.Create("bankinfo/data.go")
	if err != nil {
		log.Fatalf("Cannot create file: %q", err)
		os.Exit(1)
	}
	_, err = io.Copy(goFile, data)
	if err != nil {
		log.Fatalf("Error while writing file: %q", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func writeDataToGoFile(data []bankinfo.BankInfo) (io.Reader, error) {
	t, err := template.New("bank_data").Parse(dataTemplate)
	if err != nil {
		return nil, errors.WithMessage(err, "error while parsing template")
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, errors.WithMessage(err, "error while executing template")
	}
	formattedBytes, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, errors.WithMessage(err, "error while formatting source file")
	}
	return bytes.NewReader(formattedBytes), nil
}

const dataTemplate = `package bankinfo

var data = []BankInfo{
	{{range $element := .}}BankInfo{
		BankID: "{{.BankID}}",
		VersionNumber: "{{.VersionNumber}}",
		URL: "{{.URL}}",
		VersionName: "{{.VersionName}}",
		Institute: {{.Institute | printf "%q" }},
		City: "{{.City}}",
	},
	{{end}}
}`
