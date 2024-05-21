package model

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

var wantContactTagGroupXmlRaw string = strings.TrimSpace(`
<TagGroup>
  <GroupName>ContactInfo</GroupName>
  <Tag>
    <Name>Phone</Name>
    <Value>1-647-676-9898</Value>
  </Tag>
  <Tag>
    <Name>Email</Name>
    <Value>george@golang.com</Value>
  </Tag>
</TagGroup>
`)

var wantProductTagGroupXmlRaw string = strings.TrimSpace(`
<TagGroup>
  <GroupName>Product</GroupName>
  <Tag>
    <Name>ProductName</Name>
    <Value>Business Everyday Checking Account</Value>
  </Tag>
  <Tag>
    <Name>AccountNumber</Name>
    <Value>1234567890</Value>
  </Tag>
</TagGroup>
`)

func TestTag(t *testing.T) {
	ingestSource := &Tag{
		Name:  "ingestedSource",
		Value: "s3://bucket-name/folder-name/",
	}
	out, _ := xml.MarshalIndent(ingestSource, "", "  ")
	got := string(out)
	fmt.Println("actual: ", got)
	want := strings.TrimSpace(`
<Tag>
  <Name>ingestedSource</Name>
  <Value>s3://bucket-name/folder-name/</Value>
</Tag>
`)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProductTagGroup(t *testing.T) {
	tagGroup := makeProductTagGroup()
	out, _ := xml.MarshalIndent(tagGroup, "", "  ")
	got := string(out)
	fmt.Println("actual: ", got)
	if got != wantProductTagGroupXmlRaw {
		t.Errorf("got %q, want %q", got, wantProductTagGroupXmlRaw)
	}
}

func TestContactTagGroup(t *testing.T) {
	tagGroup := makeContactTagGroup()
	out, _ := xml.MarshalIndent(tagGroup, "", "  ")
	got := string(out)
	fmt.Println("actual: ", got)
	if got != wantContactTagGroupXmlRaw {
		t.Errorf("got %q, want %q", got, wantProductTagGroupXmlRaw)
	}
}

func TestRequestWithEmptytagsAndGroups(t *testing.T) {
	req := minRequest()
	out, _ := xml.MarshalIndent(req, "", "  ")
	got := string(out)
	want := strings.TrimSpace(`
<Request ID="1" FileName="george-2023-financial-summary.pdf" MimeType="application/pdf"></Request>
`)
	fmt.Println("actual: ", got)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRequestWithEmptyGroups(t *testing.T) {
	req := minRequest()
	// golang let you to simplify to: req.DocName = "...."
	// compiler does the heavy-lift work for you automatically
	(*req).DocName = "Customer Financial Report"
	nameValues, _ := makeTags([][]string{
		{"FirstName", "LastName", "DateOfBirth"},
		{"George", "Zhou", "1985-08-18"},
	})
	req.Metadata = &Metadata{
		Tags: *nameValues,
	}
	out, _ := xml.MarshalIndent(req, "", "  ")
	got := string(out)
	want := strings.TrimSpace(`
<Request ID="1" FileName="george-2023-financial-summary.pdf" MimeType="application/pdf" DocName="Customer Financial Report">
  <Metadata>
    <Tag>
      <Name>FirstName</Name>
      <Value>George</Value>
    </Tag>
    <Tag>
      <Name>LastName</Name>
      <Value>Zhou</Value>
    </Tag>
    <Tag>
      <Name>DateOfBirth</Name>
      <Value>1985-08-18</Value>
    </Tag>
  </Metadata>
</Request>
`)
	fmt.Println("actual: ", got)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFullRequest(t *testing.T) {
	req := minRequest()
	req.DocName = "Customer Financial Report"
	nameValues, _ := makeTags([][]string{
		{"FirstName", "LastName", "DateOfBirth"},
		{"George", "Zhou", "1985-08-18"},
	})
	req.Metadata = &Metadata{
		Tags: *nameValues,
		TagGroups: []TagGroup{
			*makeProductTagGroup(), *makeContactTagGroup(),
		},
	}
	out, _ := xml.MarshalIndent(req, "", "  ")
	got := string(out)
	want := strings.TrimSpace(`
<Request ID="1" FileName="george-2023-financial-summary.pdf" MimeType="application/pdf" DocName="Customer Financial Report">
  <Metadata>
    <Tag>
      <Name>FirstName</Name>
      <Value>George</Value>
    </Tag>
    <Tag>
      <Name>LastName</Name>
      <Value>Zhou</Value>
    </Tag>
    <Tag>
      <Name>DateOfBirth</Name>
      <Value>1985-08-18</Value>
    </Tag>
    <TagGroup>
      <GroupName>Product</GroupName>
      <Tag>
        <Name>ProductName</Name>
        <Value>Business Everyday Checking Account</Value>
      </Tag>
      <Tag>
        <Name>AccountNumber</Name>
        <Value>1234567890</Value>
      </Tag>
    </TagGroup>
    <TagGroup>
      <GroupName>ContactInfo</GroupName>
      <Tag>
        <Name>Phone</Name>
        <Value>1-647-676-9898</Value>
      </Tag>
      <Tag>
        <Name>Email</Name>
        <Value>george@golang.com</Value>
      </Tag>
    </TagGroup>
  </Metadata>
</Request>
`)
	fmt.Println("actual: ", got)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	req2 := &Request{}
	err := xml.Unmarshal([]byte(got), req2)
	if err != nil {
		t.Errorf("fail to unmarsh full request %q", err)
	}
	out2, _ := xml.MarshalIndent(req2, "", "  ")
	got2 := string(out2)
	if got != got2 {
		t.Errorf("got %q, want %q", got2, got)
	}
}

func TestRequestWithEmptytags(t *testing.T) {
	req := minRequest()
	req.DocName = "Customer Financial Report"
	req.Metadata = &Metadata{
		TagGroups: []TagGroup{
			*makeProductTagGroup(),
		},
	}
	out, _ := xml.MarshalIndent(req, "", "  ")
	got := string(out)
	want := strings.TrimSpace(`
<Request ID="1" FileName="george-2023-financial-summary.pdf" MimeType="application/pdf" DocName="Customer Financial Report">
  <Metadata>
    <TagGroup>
      <GroupName>Product</GroupName>
      <Tag>
        <Name>ProductName</Name>
        <Value>Business Everyday Checking Account</Value>
      </Tag>
      <Tag>
        <Name>AccountNumber</Name>
        <Value>1234567890</Value>
      </Tag>
    </TagGroup>
  </Metadata>
</Request>
`)
	fmt.Println("actual: ", got)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPkg01(t *testing.T) {
	req1 := makeRequest(100, "1",
		"David-01.pdf",
		"application/pdf",
		"Passport")
	david, _ := makeTags([][]string{
		{"FirstName", "LastName", "DOB"},
		{"David", "Smith", "1980-01-01"},
	})
	passport, _ := makeTags([][]string{
		{"IssueDate", "ExpiryDate", "IssuePlace"},
		{"2018-01-01", "2028-01-01", "London"},
	})
	req1.Metadata = &Metadata{
		Tags: *david,
		TagGroups: []TagGroup{
			{GroupName: "IssueInfo", Tags: *passport},
		},
	}
	req2 := makeRequest(101, "2",
		"Linda-02.pdf",
		"application/pdf",
		"Driver License")
	linda, _ := makeTags([][]string{
		{"FirstName", "LastName", "DOB"},
		{"Linda", "Chau", "1986-05-18"},
	})
	driverLicense, _ := makeTags([][]string{
		{"Grade", "ExpiryDate", "IssuePlace"},
		{"G", "2026-01-01", "Toronto"},
	})
	req2.Metadata = &Metadata{
		Tags: *linda,
		TagGroups: []TagGroup{
			{GroupName: "IssueInfo", Tags: *driverLicense},
		},
	}
	pkg := &Pkg{
		ID: "1",
		Header: PkgHeader{
			SubmissionDate: "2023-12-18",
			SubmissionTime: "15:38:18",
			Source:         "George Unit Test",
		},
		Requests: []Request{
			*req1,
			*req2,
		},
		Footer: PkgFooter{
			RequestCount: 2,
		},
	}
	out, _ := xml.MarshalIndent(pkg, "", "    ")
	got := string(out)
	want, _ := readFromFile("../testdata/model/pkg-01.xml")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func readFromFile(filename string) (string, error) {
	pwd, _ := os.Getwd()
	fmt.Println("pwd: ", pwd)
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func minRequest() *Request {
	req := &Request{
		RowNumber: 100,
		ID:        "1",
		FileName:  "george-2023-financial-summary.pdf",
		MimeType:  "application/pdf",
	}
	return req
}

func makeRequest(rowNumber int, id string, fileName string, mimeType string, docName string) *Request {
	req := &Request{
		RowNumber: rowNumber,
		ID:        id,
		FileName:  fileName,
		MimeType:  mimeType,
		DocName:   docName,
	}
	return req
}

func makeProductTagGroup() *TagGroup {
	productName := &Tag{
		Name:  "ProductName",
		Value: "Business Everyday Checking Account",
	}
	acctNumber := &Tag{
		Name:  "AccountNumber",
		Value: "1234567890",
	}
	tagGroup := &TagGroup{
		GroupName: "Product",
		Tags:      []Tag{*productName, *acctNumber},
	}
	return tagGroup
}

func makeTags(nameValuePairs [][]string) (tags *[]Tag, e error) {
	tags = nil
	e = nil
	if len(nameValuePairs) != 2 || len(nameValuePairs[0]) != len(nameValuePairs[1]) {
		e = errors.New("name value must be paired and have the same length")
		return
	}
	ret := make([]Tag, len(nameValuePairs[0]))
	for idx, name := range nameValuePairs[0] {
		ret[idx] = Tag{
			Name:  name,
			Value: nameValuePairs[1][idx],
		}
	}
	tags = &ret
	return
}

func makeContactTagGroup() *TagGroup {
	contacttags, _ := makeTags([][]string{
		{"Phone", "Email"},
		{"1-647-676-9898", "george@golang.com"},
	})
	tagGroup := &TagGroup{
		GroupName: "ContactInfo",
		Tags:      *contacttags,
	}
	return tagGroup
}
