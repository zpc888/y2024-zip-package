package service

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
	"zip-pkg-in-go/model"
)

func TestParsePackageExcelSheet1(t *testing.T) {
	testSheetWith(t, "Sheet1", "[", "]")
}

func TestParsePackageExcelSheet2(t *testing.T) {
	testSheetWith(t, "Sheet2", "", ".")
}

func testSheetWith(t *testing.T, sheetName string, groupPrefix string, groupSuffix string) {
	pi := NewParseInstruction()
	pi.SheetName = sheetName
	pi.SetGroupNameDelimiter(groupPrefix, groupSuffix)
	pkg, err := pi.ParsePackageRequests("../testdata/excel/pkg-test.xlsx")
	if err != nil || pkg == nil {
		t.Errorf("ParsePackageExcel failed: %v", err)
		return
	}
	pkg.ID = "123"
	pkg.Header = model.PkgHeader{
		SubmissionDate: "2020-01-01",
		SubmissionTime: "12:00:00",
		Source:         "UnitTest",
	}
	out, _ := xml.MarshalIndent(pkg, "", "    ")
	got := string(out)
	fmt.Println(got)
	want, _ := readFromFile("../testdata/excel/pkg-test-xlsx-expected.xml")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func readFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
