package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/mitch000001/go-hbci/generator"
)

var segmentFlag = flag.String("segment", "", "-segment 'MyAwesomeSegment'")

func main() {
	flag.Parse()
	if *segmentFlag == "" {
		fmt.Printf("You must provide a segment to generate the unmarshaler\n")
		os.Exit(1)
	}
	filename := os.Getenv("GOFILE")
	packageName := os.Getenv("GOPACKAGE")
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, filename, nil, 0)
	if err != nil {
		fmt.Println(err)
	}
	segmentGenerator := generator.NewSegmentUnmarshaler(*segmentFlag, packageName, fileSet, f)
	generated, err := segmentGenerator.Generate()
	if err != nil {
		fmt.Printf("Error while generating Unmarshaler: %v\n", err)
		os.Exit(1)
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, generated)
	if err != nil {
		fmt.Printf("Error while copying Unmarshaler: %T:%v\n", err, err)
		os.Exit(1)
	}
	newFileName := strings.TrimSuffix(filename, ".go") + "_unmarshaler.go"
	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Printf("Error while creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	fileSet = token.NewFileSet()
	newAstFile, err := parser.ParseFile(fileSet, newFileName, buf.String(), 0)
	if err != nil {
		fmt.Println(err)
	}
	err = printer.Fprint(file, fileSet, newAstFile)
	if err != nil {
		fmt.Println(err)
	}
}
