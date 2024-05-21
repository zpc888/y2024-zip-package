package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"zip-pkg-in-go/model"
)

type ParseInstruction struct {
	continuousEmptyColLimit int8
	continuousEmptyRowLimit int8
	groupPrefix             string
	groupSuffix             string
	groupIdNameDelimiter    string
	defaultMimeType         string
	defaultDocName          string
	SheetName               string
}

type colHeader struct {
	index     int
	rawName   string
	kind      int8 // 1: skip, 2: filename, 3: mimetype, 4: docname 5: id 10: tag 20: group tag
	groupId   string
	groupName string
	tagName   string
}

func (pi *ParseInstruction) parseColHeader(index int, value string) *colHeader {
	if value == "" {
		return nil
	}
	ret := &colHeader{
		index:   index,
		rawName: value,
	}
	ret.tagName = value
	if value == "Skip" {
		ret.kind = 1
	} else if value == "FileName" {
		ret.kind = 2
	} else if value == "MimeType" {
		ret.kind = 3
	} else if value == "DocName" {
		ret.kind = 4
	} else if value == "RefID" {
		ret.kind = 5
	} else {
		ret.kind = 10 // tag
		if strings.HasPrefix(value, pi.groupPrefix) {
			noPrefix := value[len(pi.groupPrefix):]
			idx := strings.Index(noPrefix, pi.groupSuffix)
			if idx != -1 {
				group := strings.TrimSpace(noPrefix[:idx])
				tag := strings.TrimSpace(noPrefix[idx+len(pi.groupSuffix):])
				if group != "" && tag != "" {
					idx2 := strings.Index(group, pi.groupIdNameDelimiter)
					if idx2 != -1 && idx2 != 0 && idx2 != len(group)-1 {
						ret.groupId = strings.TrimSpace(group[:idx2])
						ret.groupName = strings.TrimSpace(group[idx2+len(pi.groupIdNameDelimiter):])
					} else {
						ret.groupName = group
					}
					ret.tagName = tag
					ret.kind = 20 // group tag
				}
			}
		}
	}
	return ret
}

func NewParseInstruction() *ParseInstruction {
	return &ParseInstruction{
		continuousEmptyColLimit: 10,
		continuousEmptyRowLimit: 10,
		groupPrefix:             "[",
		groupSuffix:             "]",
		groupIdNameDelimiter:    ":",
		SheetName:               "Sheet1",
	}
}

func (pi *ParseInstruction) SetContinuousEmptyColLimit(limit int8) {
	pi.continuousEmptyColLimit = limit
}

func (pi *ParseInstruction) SetContinuousEmptyRowLimit(limit int8) {
	pi.continuousEmptyRowLimit = limit
}

func (pi *ParseInstruction) SetGroupNameDelimiter(prefix string, suffix string) {
	pi.groupPrefix = prefix
	pi.groupSuffix = suffix
}

func (pi *ParseInstruction) SetGroupIdNameDelimiter(delim string) {
	pi.groupIdNameDelimiter = delim
}

func (pi *ParseInstruction) ParsePackageRequests(xlsx string) (*model.Pkg, error) {
	f, err := excelize.OpenFile(xlsx)
	if err != nil {
		fmt.Printf("Error openinng excel file: %v\n", err)
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing excel file: %v\n", err)
		}
	}()
	rows, err := f.GetRows(pi.SheetName)
	if err != nil {
		fmt.Printf("Error get excel rows: %v\n", err)
		return nil, err
	}
	ret := &model.Pkg{}
	ret.Requests = make([]model.Request, 0)

	var headerMap map[int]*colHeader = make(map[int]*colHeader)
	var continueEmptyRowCount int8 = 0
	var maxColIdx = 0
	var seq = 0
	for i, row := range rows {
		if i == 0 { // header row
			maxColIdx = pi.parseHeaderRow(row, headerMap)
		} else {
			req, status := pi.buildRequestAndStatus(row, &headerMap, maxColIdx)
			if status == 2 { // empty row
				continueEmptyRowCount += 1
				if continueEmptyRowCount > pi.continuousEmptyRowLimit {
					break
				}
			} else if status == 1 { // ignore row
				continueEmptyRowCount = 0
			} else {
				continueEmptyRowCount = 0
				req.RowNumber = i
				seq += 1
				if req.ID == "" {
					req.ID = strconv.Itoa(seq)
				}
				if req.DocName == "" && pi.defaultDocName != "" {
					req.DocName = pi.defaultDocName
				}
				if req.MimeType == "" && pi.defaultMimeType != "" {
					req.MimeType = pi.defaultMimeType
				}
				ret.Requests = append(ret.Requests, *req)
			}
		}
	}
	ret.Footer = model.PkgFooter{
		RequestCount: int16(len(ret.Requests)),
	}
	return ret, nil
}

func (pi *ParseInstruction) buildRequestAndStatus(row []string, headerMap *map[int]*colHeader, maxColIdx int) (*model.Request, int8) {
	req := &model.Request{
		Metadata: &model.Metadata{},
	}
	var status int8 = 2 // 0: ok; 1: ignore; 2: all empty
	for j, cell := range row {
		if j > maxColIdx {
			break
		}
		col := strings.TrimSpace(cell)
		if col == "" {
			continue
		}
		if header, ok := (*headerMap)[j]; ok {
			status = 0
			if header.kind == 1 && (strings.EqualFold(col, "yes") || strings.EqualFold(col, "true")) {
				status = 1
				break
			}
			if header.kind == 2 {
				req.FileName = col
			} else if header.kind == 3 {
				req.MimeType = col
			} else if header.kind == 4 {
				req.DocName = col
			} else if header.kind == 5 {
				req.ID = col
			} else {
				req.Metadata.AddTagOrGroupTag(header.groupId, header.groupName, header.tagName, col)
			}
		}
	}
	return req, status
}

func (pi *ParseInstruction) parseHeaderRow(row []string, headerMap map[int]*colHeader) int {
	var maxColIdx int = -1
	var continueEmptyColCount int8 = 0
	for j, cell := range row {
		header := pi.parseColHeader(j, strings.TrimSpace(cell))
		if header != nil {
			headerMap[j] = header
			maxColIdx = j
			continueEmptyColCount = 0
		} else {
			continueEmptyColCount += 1
			if continueEmptyColCount > pi.continuousEmptyColLimit {
				break
			}
		}
	}
	return maxColIdx
}
