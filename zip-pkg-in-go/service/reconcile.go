package service

import (
	"encoding/xml"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"zip-pkg-in-go/model"
)

type ReconcileInstruction struct {
	ReportDir          string
	OutDir             string
	ReportFileEndsWith string
}

type ReconcileResult struct {
	Request  *model.Request
	Document *model.ReportDocument
}

func NewReconcileInstruction() *ReconcileInstruction {
	return &ReconcileInstruction{
		ReportDir:          "report",
		OutDir:             "output",
		ReportFileEndsWith: ".xml",
	}
}

func (ri *ReconcileInstruction) Reconcile(pkg *model.Pkg) (*[]ReconcileResult, error) {
	reportDocs, err := ri.readAllReports()
	if err != nil {
		return nil, err
	}
	var join = make(map[string]ReconcileResult)
	for idx, doc := range *reportDocs {
		// &doc will be the same always because it's a pointer to the same object, but its value will be different for each loop
		fmt.Printf("Reconcile report doc - File Name = %v; doc address = %p vs %p\n", doc.FileName, &doc, &((*reportDocs)[idx]))
		join[doc.FileName] = ReconcileResult{
			Request:  nil,
			Document: &((*reportDocs)[idx]),
		}
	}
	var tmpResults = make([]ReconcileResult, 0)
	inBothOk, inBothErr, inReqOnly, inRepOk, inRepErr := 0, 0, 0, 0, 0
	for idx, req := range pkg.Requests {
		// &req will be the same always because it's a pointer to the same object, but its value will be different for each loop
		fmt.Printf("Reconcile request #%v - File Name = %v; req address = %p vs %p\n", idx+1, req.FileName, &req, &(pkg.Requests[idx]))
		if existing, ok := join[req.FileName]; !ok {
			tmpResults = append(tmpResults, ReconcileResult{
				Request:  &(pkg.Requests[idx]),
				Document: nil,
			})
			inReqOnly++
		} else {
			existing.Request = &(pkg.Requests[idx])
			tmpResults = append(tmpResults, existing)
			delete(join, req.FileName)
			if existing.Document.Status == "Succeeded" {
				inBothOk++
			} else {
				inBothErr++
			}
		}
	}
	for _, tmp := range join {
		tmpResults = append(tmpResults, tmp)
		if tmp.Document.Status == "Succeeded" {
			inRepOk++
		} else {
			inRepErr++
		}
	}
	fmt.Println()
	fmt.Println("Reconcile result:")
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Requests OK           : %v\n", inBothOk)
	fmt.Printf("Requests Error        : %v\n", inBothErr)
	fmt.Printf("Requests (No Response): %v\n", inReqOnly)
	fmt.Printf("Responses OK (No Req) : %v\n", inRepOk)
	fmt.Printf("Responses Err(No Req) : %v\n", inRepErr)
	fmt.Printf("Total %v requests vs %v responses\n", inBothOk+inBothErr+inReqOnly, inBothOk+inBothErr+inRepOk+inRepErr)
	fmt.Println("----------------------------------------------------------------")
	fmt.Println()
	return &tmpResults, nil
}

func (ri *ReconcileInstruction) readAllReports() (*[]model.ReportDocument, error) {
	filenames := make([]string, 0)
	filepath.Walk(ri.ReportDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ri.ReportFileEndsWith) {
			filenames = append(filenames, path)
		}
		return nil
	})
	var docs = make([]model.ReportDocument, 0)
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		report := &model.Report{}
		_ = xml.Unmarshal(data, report)
		docs = append(docs, report.Documents...)
	}
	return &docs, nil
}

func (ri *ReconcileInstruction) OutputExcel(results *[]ReconcileResult, requestHeaders *[]ColHeader, excel *excelize.File) error {
	index, _ := excel.NewSheet("Reconcile")
	colMap := colIndexToString(len(*requestHeaders))
	rowNum := 1
	r := strconv.Itoa(rowNum)
	for i, header := range *requestHeaders {
		err := excel.SetCellValue("Reconcile", colMap[i]+r, header.RawName)
		if err != nil {
			return err
		}
	}
	repIdxFrom := len(*requestHeaders)
	_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom]+r, "Report Status")
	_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom+1]+r, "Content ID")
	_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom+2]+r, "Error Code")
	_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom+3]+r, "Error Message")
	_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom+4]+r, "Report DocID")
	for _, result := range *results {
		rowNum++
		r = strconv.Itoa(rowNum)
		ri.outputResultRow(result.Request, requestHeaders, repIdxFrom, result.Document, colMap, excel, r)
	}
	excel.SetActiveSheet(index)
	return nil
}

func (ri *ReconcileInstruction) outputResultRow(req *model.Request, reqHeaders *[]ColHeader, repIdxFrom int, doc *model.ReportDocument, colMap map[int]string, excel *excelize.File, rowNum string) {
	if req != nil {
		for i := 0; i < repIdxFrom; i++ {
			colHeader := (*reqHeaders)[i]
			val := resolveRequestColValue(req, colHeader)
			_ = excel.SetCellValue("Reconcile", colMap[i]+rowNum, val)
		}
	}
	if doc != nil {
		_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom]+rowNum, doc.Status)
		_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom+1]+rowNum, doc.ContentID)
		_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom+2]+rowNum, doc.ErrorCode)
		_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom+3]+rowNum, doc.ErrorMessage)
		_ = excel.SetCellValue("Reconcile", colMap[repIdxFrom+4]+rowNum, doc.ID+" - "+doc.FileName)
	}
}

func resolveRequestColValue(req *model.Request, header ColHeader) string {
	if header.Kind == 1 { // skip
		return ""
	} else if header.Kind == 2 { // filename
		return req.FileName
	} else if header.Kind == 3 {
		return req.MimeType
	} else if header.Kind == 4 {
		return req.DocName
	} else if header.Kind == 5 {
		return req.ID
	} else if header.Kind == 10 { // tag
		return req.GetTagValue(header.TagName)
	} else if header.Kind == 20 { // group tag
		return req.GetTagGroupValue(header.GroupId, header.GroupName, header.TagName)
	} else {
		return ""
	}
}

func colIndexToString(colIndex int) map[int]string {
	f := func(idx int) string {
		ret := ""
		idx2 := idx
		for {
			ret = string(rune('A'+(idx2%26))) + ret
			idx2 = idx2 / 26
			if idx2 == 0 {
				break
			}
			idx2--
		}
		return ret
	}
	ret := make(map[int]string)
	for i := 0; i < colIndex; i++ {
		ret[i] = f(i)
	}
	// status, content id, error code, error message, report docID
	for i := colIndex; i < colIndex+5; i++ {
		ret[i] = f(i)
	}
	return ret
}
