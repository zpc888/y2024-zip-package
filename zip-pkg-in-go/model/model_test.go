package model

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"testing"
)

var wantContactAttrGroupXmlRaw string = strings.TrimSpace(`
<AttributeGroup>
  <GroupName>ContactInfo</GroupName>
  <Attribute>
    <Name>Phone</Name>
    <Value>1-647-676-9898</Value>
  </Attribute>
  <Attribute>
    <Name>Email</Name>
    <Value>george@golang.com</Value>
  </Attribute>
</AttributeGroup>
`)

var wantProductAttrGroupXmlRaw string = strings.TrimSpace(`
<AttributeGroup>
  <GroupName>Product</GroupName>
  <Attribute>
    <Name>ProductName</Name>
    <Value>Business Everyday Checking Account</Value>
  </Attribute>
  <Attribute>
    <Name>AccountNumber</Name>
    <Value>1234567890</Value>
  </Attribute>
</AttributeGroup>
`)

func TestAttribute(t *testing.T) {
	ingestSource := &Attribute{
		Name:  "ingestedSource",
		Value: "s3://bucket-name/folder-name/",
	}
	out, _ := xml.MarshalIndent(ingestSource, "", "  ")
	got := string(out)
	fmt.Println("actual: ", got)
	want := strings.TrimSpace(`
<Attribute>
  <Name>ingestedSource</Name>
  <Value>s3://bucket-name/folder-name/</Value>
</Attribute>
`)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProductAttributeGroup(t *testing.T) {
	attributeGroup := makeProductAttrGroup()
	out, _ := xml.MarshalIndent(attributeGroup, "", "  ")
	got := string(out)
	fmt.Println("actual: ", got)
	if got != wantProductAttrGroupXmlRaw {
		t.Errorf("got %q, want %q", got, wantProductAttrGroupXmlRaw)
	}
}

func TestContactAttributeGroup(t *testing.T) {
	attributeGroup := makeContactAttrGroup()
	out, _ := xml.MarshalIndent(attributeGroup, "", "  ")
	got := string(out)
	fmt.Println("actual: ", got)
	if got != wantContactAttrGroupXmlRaw {
		t.Errorf("got %q, want %q", got, wantProductAttrGroupXmlRaw)
	}
}

func TestRequestWithEmptyAttrsAndGroups(t *testing.T) {
	req := minRequest()
	out, _ := xml.MarshalIndent(req, "", "  ")
	got := string(out)
	want := strings.TrimSpace(`
<Request RefID="1" FileName="george-2023-financial-summary.pdf" MimeType="application/pdf"></Request>
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
	nameValues, _ := makeAttributes([][]string{
		{"FirstName", "LastName", "DateOfBirth"},
		{"George", "Zhou", "1985-08-18"},
	})
	req.Metadata = &Metadata{
		Attributes: *nameValues,
	}
	out, _ := xml.MarshalIndent(req, "", "  ")
	got := string(out)
	want := strings.TrimSpace(`
<Request RefID="1" FileName="george-2023-financial-summary.pdf" MimeType="application/pdf" DocName="Customer Financial Report">
  <Metadata>
    <Attribute>
      <Name>FirstName</Name>
      <Value>George</Value>
    </Attribute>
    <Attribute>
      <Name>LastName</Name>
      <Value>Zhou</Value>
    </Attribute>
    <Attribute>
      <Name>DateOfBirth</Name>
      <Value>1985-08-18</Value>
    </Attribute>
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
	nameValues, _ := makeAttributes([][]string{
		{"FirstName", "LastName", "DateOfBirth"},
		{"George", "Zhou", "1985-08-18"},
	})
	req.Metadata = &Metadata{
		Attributes: *nameValues,
		AttributeGroups: []AttributeGroup{
			*makeProductAttrGroup(), *makeContactAttrGroup(),
		},
	}
	out, _ := xml.MarshalIndent(req, "", "  ")
	got := string(out)
	want := strings.TrimSpace(`
<Request RefID="1" FileName="george-2023-financial-summary.pdf" MimeType="application/pdf" DocName="Customer Financial Report">
  <Metadata>
    <Attribute>
      <Name>FirstName</Name>
      <Value>George</Value>
    </Attribute>
    <Attribute>
      <Name>LastName</Name>
      <Value>Zhou</Value>
    </Attribute>
    <Attribute>
      <Name>DateOfBirth</Name>
      <Value>1985-08-18</Value>
    </Attribute>
    <AttributeGroup>
      <GroupName>Product</GroupName>
      <Attribute>
        <Name>ProductName</Name>
        <Value>Business Everyday Checking Account</Value>
      </Attribute>
      <Attribute>
        <Name>AccountNumber</Name>
        <Value>1234567890</Value>
      </Attribute>
    </AttributeGroup>
    <AttributeGroup>
      <GroupName>ContactInfo</GroupName>
      <Attribute>
        <Name>Phone</Name>
        <Value>1-647-676-9898</Value>
      </Attribute>
      <Attribute>
        <Name>Email</Name>
        <Value>george@golang.com</Value>
      </Attribute>
    </AttributeGroup>
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

func TestRequestWithEmptyAttrs(t *testing.T) {
	req := minRequest()
	req.DocName = "Customer Financial Report"
	req.Metadata = &Metadata{
		AttributeGroups: []AttributeGroup{
			*makeProductAttrGroup(),
		},
	}
	out, _ := xml.MarshalIndent(req, "", "  ")
	got := string(out)
	want := strings.TrimSpace(`
<Request RefID="1" FileName="george-2023-financial-summary.pdf" MimeType="application/pdf" DocName="Customer Financial Report">
  <Metadata>
    <AttributeGroup>
      <GroupName>Product</GroupName>
      <Attribute>
        <Name>ProductName</Name>
        <Value>Business Everyday Checking Account</Value>
      </Attribute>
      <Attribute>
        <Name>AccountNumber</Name>
        <Value>1234567890</Value>
      </Attribute>
    </AttributeGroup>
  </Metadata>
</Request>
`)
	fmt.Println("actual: ", got)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func minRequest() *Request {
	req := &Request{
		RowNumber: 100,
		RefID:     "1",
		FileName:  "george-2023-financial-summary.pdf",
		MimeType:  "application/pdf",
	}
	return req
}

func makeProductAttrGroup() *AttributeGroup {
	productName := &Attribute{
		Name:  "ProductName",
		Value: "Business Everyday Checking Account",
	}
	acctNumber := &Attribute{
		Name:  "AccountNumber",
		Value: "1234567890",
	}
	attributeGroup := &AttributeGroup{
		GroupName:  "Product",
		Attributes: []Attribute{*productName, *acctNumber},
	}
	return attributeGroup
}

func makeAttributes(nameValuePairs [][]string) (attrs *[]Attribute, e error) {
	attrs = nil
	e = nil
	if len(nameValuePairs) != 2 || len(nameValuePairs[0]) != len(nameValuePairs[1]) {
		e = errors.New("name value must be paired and have the same length")
		return
	}
	ret := make([]Attribute, len(nameValuePairs[0]))
	for idx, name := range nameValuePairs[0] {
		ret[idx] = Attribute{
			Name:  name,
			Value: nameValuePairs[1][idx],
		}
	}
	attrs = &ret
	return
}

func makeContactAttrGroup() *AttributeGroup {
	contactAttrs, _ := makeAttributes([][]string{
		{"Phone", "Email"},
		{"1-647-676-9898", "george@golang.com"},
	})
	attributeGroup := &AttributeGroup{
		GroupName:  "ContactInfo",
		Attributes: *contactAttrs,
	}
	return attributeGroup
}
