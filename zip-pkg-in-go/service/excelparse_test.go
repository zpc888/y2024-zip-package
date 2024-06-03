package service

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
	"zip-pkg-in-go/model"
)

func TestParsePackageExcelSheet1(t *testing.T) {
	testSheetWith(t, "Sheet1", "[", "]", nil)
}

func TestParsePackageExcelSheet2(t *testing.T) {
	testSheetWith(t, "Sheet2", "", ".", func(pkg *model.Pkg) {
		req := pkg.Requests[0]
		req.FileName = "David2-Passport.pdf"
		req.ID = "3"

		// ERROR #1: The below code is not working, since req.Metadata is a pointer, so changing it will affect both 1st and 3rd request in output
		// req.Metadata.Tags[0].Value = "David2"

		// ERROR #2: The below code is not working, since tag is not a pointer, but a value struct. Changing tag.value won't affect req.Metadata.Tags
		//for _, tag := range req.Metadata.Tags {
		//	if tag.Name == "FirstName" {
		//		tag.Value = "David2"
		//	}
		//}

		// Correct way to do -- STEP 1
		md3 := *(req.Metadata) // copy the metadata into different variable
		fmt.Printf("Request #1 Metadata pointer address: %p\n", req.Metadata)
		fmt.Printf("Request #3 Metadata pointer address: %p\n", &md3)
		fmt.Printf("Request #1 Metadata.Tags pointer address: %p\n", &(req.Metadata.Tags))
		fmt.Printf("Request #3 Metadata.Tags pointer address: %p\n", &(md3.Tags))
		fmt.Printf("Request #1 Metadata.Tags[0] pointer address: %p\n", &(req.Metadata.Tags[0]))
		fmt.Printf("Request #3 Metadata.Tags[0] pointer address: %p\n", &(md3.Tags[0]))
		//Request #1 Metadata pointer address: 0xc0004038c0
		//Request #3 Metadata pointer address: 0xc000403950
		//Request #1 Metadata.Tags pointer address: 0xc0004038c0
		//Request #3 Metadata.Tags pointer address: 0xc000403950
		//Request #1 Metadata.Tags[0] pointer address: 0xc000360c00
		//Request #3 Metadata.Tags[0] pointer address: 0xc000360c00

		// ERROR #3: not working since []Tag is a slice, although the address is different, but the content address is a pointer to the same memory address,
		//           i.e. req.Metadata.Tags[0] and md3.Tags[0] are pointing to the same memory address
		//           but req.Metadata.Tags and md3.Tags are pointing to different memory address
		// md3.Tags[0].Value = "David2"

		// Correct way to do -- STEP 2

		tags := make([]model.Tag, 0)
		tags = append(tags, req.Metadata.Tags...)
		tags[0].Value = "David2"
		fmt.Printf("new tags[0] pointer address: %p\n", &(tags[0]))
		md3.Tags = tags

		fmt.Printf("Later Request #1 Metadata.Tags[0] pointer address: %p\n", &(req.Metadata.Tags[0]))
		req.Metadata = &md3 // assign the copied metadata pointer back to request #3, so that request #1 won't be affected
		fmt.Printf("Later Request #3 Metadata.Tags[0] pointer address: %p\n", &(md3.Tags[0]))

		pkg.Requests = append(pkg.Requests, req)
		pkg.Trailer.RequestCount = pkg.Trailer.RequestCount + 1
	})
}

func testSheetWith(t *testing.T, sheetName string, groupPrefix string, groupSuffix string, pkgRefiner func(pkg *model.Pkg)) {
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
	pkg.Trailer = model.PkgTrailer{
		RequestCount: int16(len(pkg.Requests)),
	}
	out, _ := xml.MarshalIndent(pkg, "", "    ")
	got := string(out)
	fmt.Println("Got", got)
	want, _ := readFromFile("../testdata/excel/pkg-test-xlsx-expected.xml")
	expectedPkg := model.Pkg{}
	_ = xml.Unmarshal([]byte(want), &expectedPkg)
	if pkgRefiner != nil {
		pkgRefiner(&expectedPkg)
	}
	expectedOut, _ := xml.MarshalIndent(&expectedPkg, "", "    ")
	want = string(expectedOut)
	fmt.Println("Want", want)
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
