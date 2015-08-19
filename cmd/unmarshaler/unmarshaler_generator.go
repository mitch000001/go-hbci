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
	"strconv"
	"strings"

	"github.com/mitch000001/go-hbci/generator"
)

var segmentFlag = flag.String("segment", "", "'MyAwesomeSegment'")
var segmentInterfaceFlag = flag.String("segment_interface", "Segment", "'MyAwesomeInterface'")
var segmentVersionsFlag segmentVersions

func init() {
	flag.Var(&segmentVersionsFlag, "segment_versions", "'MyAwesomeSegmentVersion1:1,MyAwesomeSegmentVersion2:2'")
}

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
	segment := generator.SegmentIdentifier{
		Name:          *segmentFlag,
		InterfaceName: *segmentInterfaceFlag,
		Versions:      segmentVersionsFlag,
	}
	var generated io.Reader
	if len(segmentVersionsFlag) != 0 {
		segmentGenerator := generator.NewVersionedSegmentUnmarshaler(segment, packageName, fileSet, f)
		generated, err = segmentGenerator.Generate()
	} else {
		segmentGenerator := generator.NewSegmentUnmarshaler(segment, packageName, fileSet, f)
		generated, err = segmentGenerator.Generate()
	}
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

type segmentVersions []generator.SegmentIdentifier

func (s *segmentVersions) String() string {
	var buf bytes.Buffer
	for _, version := range *s {
		fmt.Fprintf(&buf, "%s:%d", version.Name, version.Version)
	}
	return buf.String()
}

func (s *segmentVersions) Set(in string) error {
	unquoted, err := strconv.Unquote(in)
	if err != nil {
		return fmt.Errorf("Invalid input: %q (%v)", in, err)
	}
	segments := strings.FieldsFunc(unquoted, func(r rune) bool {
		return r == ','
	})
	for _, seg := range segments {
		parts := strings.SplitN(seg, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Malformed versioned segment: %q", seg)
		}
		version, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("Malformed segment version: %v", err)
		}
		*s = append(*s, generator.SegmentIdentifier{Name: parts[0], Version: version})
	}
	return nil
}
