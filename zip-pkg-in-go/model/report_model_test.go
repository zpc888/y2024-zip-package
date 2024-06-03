package model

import (
	"encoding/xml"
	"testing"
)

func TestToAndFromXml(t *testing.T) {
	r := buildTestReport()
	out, _ := xml.MarshalIndent(r, "", "  ")
	got := string(out)
	fromXml, _ := readFromFile("../testdata/excel/pkg-test-xlsx-report.xml")
	r2 := &Report{}
	_ = xml.Unmarshal([]byte(fromXml), r2)
	out2, _ := xml.MarshalIndent(r2, "", "  ")
	want := string(out2)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func buildTestReport() *Report {
	r := minReport()
	doc1 := newOkDocument("1", "David-Passport.pdf", "A73jdf83838D")
	addTags(doc1, []string{"FirstName", "David", "LastName", "Smith", "DOB", "1986-05-18"})
	addTagGroup(doc1, "IssueInfo", []string{"IssueDate", "2011/01/01", "ExpiryDate", "2021/01/01", "IssuePlace", "Toronto"})
	addTagGroup(doc1, "IssueInfo", []string{"IssuePlace", "Canada"})
	addTagGroup(doc1, "ContactInfo", []string{"Phone", "647-875-8899", "Email", "david.smith@gmail.com"})
	addDocument(r, doc1)

	doc2 := newOkDocument("2", "Linda-DriverLicense.png", "A73jdf83839D")
	addTags(doc2, []string{"FirstName", "Linda", "LastName", "Chau", "DOB", "1988/01/06"})
	addTagGroup(doc2, "IssueInfo", []string{"IssueDate", "2016/01/01", "ExpiryDate", "2026/01/01", "Grade", "G", "IssuePlace", "London"})
	addTagGroup(doc2, "IssueInfo", []string{"IssuePlace", "Canada"})
	addTagGroup(doc2, "ContactInfo", []string{"Phone", "437-441-1564", "Email", "Linda.Chau@yahoo.com"})
	addDocument(r, doc2)

	doc3 := makeErrorDocument("3", "pp-1006.pdf", "8738", "Invalid file format")
	addDocument(r, doc3)
	return r
}

func minReport() *Report {
	return &Report{
		ID: "123",
		Header: ReportHeader{
			SubmissionDate:     "2020-01-01",
			SubmissionTime:     "11:18:23",
			RequestApplication: "unit-test",
			PackageName:        "gz-unit-test-01.zip",
			ContentType:        "ABC8789",
			ProcessingDuration: "00:00:02.883",
			ProcessingDate:     "2020-08-08",
			ProcessingTime:     "12:25:18",
		},
		Documents: []ReportDocument{},
		Trailer: ReportTrailer{
			DocumentCount: 0,
			SuccessCount:  0,
			ErrorCount:    0,
		},
	}
}

func addDocument(r *Report, doc *ReportDocument) {
	r.Documents = append(r.Documents, *doc)
	r.Trailer.DocumentCount++
	if doc.Status == "Succeeded" {
		r.Trailer.SuccessCount++
	} else {
		r.Trailer.ErrorCount++
	}
}

func addTags(doc *ReportDocument, tagNameValues []string) {
	for i := 0; i < len(tagNameValues); i += 2 {
		doc.Metadata.Tags = append(doc.Metadata.Tags, Tag{
			Name:  tagNameValues[i],
			Value: tagNameValues[i+1],
		})
	}
}

func addTagGroup(doc *ReportDocument, tagGroupName string, tagNameValues []string) {
	tagGroup := TagGroup{
		GroupName: tagGroupName,
		Tags:      []Tag{},
	}
	for i := 0; i < len(tagNameValues); i += 2 {
		tagGroup.Tags = append(tagGroup.Tags, Tag{
			Name:  tagNameValues[i],
			Value: tagNameValues[i+1],
		})
	}
	doc.Metadata.TagGroups = append(doc.Metadata.TagGroups, tagGroup)
}

func newOkDocument(id, fileName, contentId string) *ReportDocument {
	return &ReportDocument{
		ID:        id,
		FileName:  fileName,
		Status:    "Succeeded",
		ContentID: contentId,
		Metadata: &Metadata{
			Tags:      []Tag{},
			TagGroups: make([]TagGroup, 0),
		},
	}
}

func makeErrorDocument(id, fileName, errorCode, errorMessage string) *ReportDocument {
	return &ReportDocument{
		ID:           id,
		FileName:     fileName,
		Status:       "Failed",
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}
