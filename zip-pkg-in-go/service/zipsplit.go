package service

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"zip-pkg-in-go/model"
)

type ZipInstruction struct {
	SrcDir                string
	DstDir                string
	MaxSize               int64
	TargetFileNamePattern string
	Unzip                 bool
	SourceID              string
	MetaXmlFileName       string
}

func NewZipInstruction() *ZipInstruction {
	return &ZipInstruction{
		SrcDir:                "sources",
		DstDir:                "output",
		MaxSize:               980 * 1024 * 1024,
		TargetFileNamePattern: "package-${yyMMddHHmmssSSS}-${splitSeq}",
		Unzip:                 true,
		SourceID:              "0086",
		MetaXmlFileName:       "package-metadata.xml",
	}
}

func (zi *ZipInstruction) Zip(requests *[]model.Request) error {
	err := ensureDir(zi.SrcDir)
	if err != nil {
		return err
	}
	size := int64(0)
	splitSeq := 1
	fromIdx := 0
	for i, req := range *requests {
		fileName := zi.SrcDir + "/" + req.FileName
		info, err2 := os.Stat(fileName)
		if err2 != nil {
			return err2
		}
		size += info.Size()
		if size > zi.MaxSize {
			zipErr := zi.zipFiles((*requests)[fromIdx:i], splitSeq)
			if zipErr != nil {
				return zipErr
			}
			splitSeq++
			fromIdx = i
			size = info.Size()
		}
	}
	zipErr := zi.zipFiles((*requests)[fromIdx:], splitSeq)
	if zipErr != nil {
		return zipErr
	}
	return nil
}

func (zi *ZipInstruction) zipFiles(requests []model.Request, seq int) error {
	err := os.MkdirAll(zi.DstDir, 0755)
	if err != nil {
		return err
	}
	fn := zi.resolveTargetFileName(seq)
	f, e := os.Create(zi.DstDir + "/" + fn + ".zip")
	if e != nil {
		return e
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	zipWriter := zip.NewWriter(f)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			panic(err)
		}
	}(zipWriter)
	unzipD := zi.DstDir + "/" + fn + ".d"
	if zi.Unzip {
		mkDirE := os.Mkdir(unzipD, 0755)
		if mkDirE != nil {
			return mkDirE
		}
	}
	for _, req := range requests {
		zipE := zi.doZipFile(zipWriter, req)
		if zipE != nil {
			return zipE
		}
		if zi.Unzip {
			cpE := zi.doCopySourceFile(req, unzipD)
			if cpE != nil {
				return cpE
			}
		}
	}
	pkg := &model.Pkg{
		ID: strconv.Itoa(seq),
		Header: model.PkgHeader{
			SubmissionDate: time.Now().Format("2006-01-02"),
			SubmissionTime: time.Now().Format("15:04:05"),
			Source:         zi.SourceID,
		},
		Trailer: model.PkgTrailer{
			RequestCount: int16(len(requests)),
		},
		Requests: requests,
	}
	xmlBytes, _ := xml.MarshalIndent(pkg, "", "    ")
	xmlContent := string(xmlBytes)
	zipE := zi.doZipMetaXml(zipWriter, &xmlContent)
	if zipE != nil {
		return zipE
	}
	if zi.Unzip {
		cpE := zi.doCopyMetaXml(&xmlContent, unzipD)
		if cpE != nil {
			return cpE
		}
	}
	return nil
}

func (zi *ZipInstruction) doZipMetaXml(zw *zip.Writer, xmlStr *string) error {
	w, we := zw.Create(zi.MetaXmlFileName)
	if we != nil {
		return we
	}
	if _, err := io.WriteString(w, *xmlStr); err != nil {
		return err
	}
	fmt.Println("zipped: ", zi.MetaXmlFileName)
	return nil
}

func (zi *ZipInstruction) doZipFile(zw *zip.Writer, req model.Request) error {
	f, e := os.Open(zi.SrcDir + "/" + req.FileName)
	if e != nil {
		return e
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	w, we := zw.Create(req.FileName)
	if we != nil {
		return we
	}
	if _, err := io.Copy(w, f); err != nil {
		return err
	}
	fmt.Println("zipped: ", req.FileName)
	return nil
}

func (zi *ZipInstruction) doCopyMetaXml(xmlStr *string, dstDir string) error {
	w, we := os.Create(dstDir + "/" + zi.MetaXmlFileName)
	if we != nil {
		return we
	}
	if _, err := io.WriteString(w, *xmlStr); err != nil {
		return err
	}
	err := w.Close()
	if err != nil {
		return err
	}
	fmt.Println("copied: ", zi.MetaXmlFileName)
	return nil
}

func (zi *ZipInstruction) doCopySourceFile(req model.Request, dstDir string) error {
	f, e := os.Open(zi.SrcDir + "/" + req.FileName)
	if e != nil {
		return e
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	w, we := os.Create(dstDir + "/" + req.FileName)
	if we != nil {
		return we
	}
	if _, err := io.Copy(w, f); err != nil {
		return err
	}
	err := w.Close()
	if err != nil {
		return err
	}
	fmt.Println("copied: ", req.FileName)
	return nil
}

func (zi *ZipInstruction) resolveTargetFileName(seq int) string {
	tm := time.Now()
	ts := fmt.Sprintf("%02d%02d%02d%02d%02d%02d%03d", tm.Year()%100, tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond()/1000000)
	fileName := strings.ReplaceAll(zi.TargetFileNamePattern, "${yyMMddHHmmssSSS}", ts)
	fileName = strings.ReplaceAll(fileName, "${splitSeq}", strconv.Itoa(seq))
	return fileName
}

func ensureDir(dir string) error {
	srcDir, err := os.Open(dir)
	if err != nil {
		return err
	}
	srcInfo, err2 := srcDir.Stat()
	if err2 != nil {
		return err2
	}
	if !srcInfo.IsDir() {
		return errors.New("source [" + dir + "] is not a directory")
	}
	return nil
}
