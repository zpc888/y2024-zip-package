package main

import (
	"bufio"
	_ "flag"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
	"strings"
	"time"
	"zip-pkg-in-go/service"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "--help" || os.Args[1] == "-h" || len(os.Args)%2 == 0 {
		usage()
		os.Exit(1)
	}
	var cmd, fileDir, outDir, reportDir, fileEndsWith, xls, config, sheetName string = "", "sources", "output", "report", ".xml", "metadata.xlsx", "config.properties", "Sheet1"
	unzip := true
	for i := 1; i < len(os.Args); i += 2 {
		switch os.Args[i] {
		case "--command":
			cmd = os.Args[i+1]
		case "--file-dir":
			fileDir = os.Args[i+1]
		case "--out-dir":
			outDir = os.Args[i+1]
		case "--report-dir":
			reportDir = os.Args[i+1]
		case "--xls":
			xls = os.Args[i+1]
		case "--config":
			config = os.Args[i+1]
		case "--report-file-ends-with":
			fileEndsWith = os.Args[i+1]
		case "--unzip":
			unzip = os.Args[i+1] != "false"
		case "--sheet-name":
			sheetName = os.Args[i+1]
		default:
			fmt.Println("Invalid option")
		}
	}
	start := time.Now()
	if cmd == "package" {
		pkg(fileDir, outDir, xls, config, sheetName, unzip)
		duration := time.Since(start)
		fmt.Printf("Duration: %v\n", duration)
	} else if cmd == "reconcile" {
		reconcile(reportDir, fileEndsWith, outDir, xls, config, sheetName)
		duration := time.Since(start)
		fmt.Printf("Duration: %v\n", duration)
	} else {
		usage()
		os.Exit(1)
	}
}

func pkg(srcDir, outDir, xls, config, sheetName string, unzip bool) {
	fmt.Printf("Package: %s %s %s %s %s\n", srcDir, outDir, xls, config, sheetName)
	cfg := loadConfig(config)
	pi := service.NewParseInstruction()
	pi.SheetName = sheetName
	configParseInstructure(pi, cfg)
	pkg, err := pi.ParsePackageRequests(xls)
	if err != nil || pkg == nil {
		fmt.Printf("ParsePackageExcel failed or no requests at all: %v\n", err)
		return
	}
	fmt.Println("Parse success and get requests: ", len(pkg.Requests))
	zi := service.NewZipInstruction()
	zi.SrcDir = srcDir
	zi.DstDir = outDir
	zi.Unzip = unzip
	configZipInstructure(zi, cfg)
	err = zi.Zip(&(pkg.Requests))
	if err != nil {
		fmt.Printf("Zip failed: %v\n", err)
		return
	}
	fmt.Println("Zip success")
}

func configZipInstructure(pi *service.ZipInstruction, cfg *map[string]string) {
	//obtain zip instruction info from config to set the following values
	//zi.MaxSize = 980 * 1024 * 1024
	//zi.TargetFileNamePattern = "package-${yyMMddHHmmssSSS}-${splitSeq}"
	//zi.SourceID = "0086"
	//zi.MetaXmlFileName = "package-metadata.xml"
	maxSize, ok10 := (*cfg)["zip-package-max-size"]
	if ok10 {
		unit := 1024 * 1024 // 1MB
		suffixLen := 0
		maxSize = strings.ToLower(maxSize)
		if strings.HasSuffix(maxSize, "gb") {
			unit = 1024 * unit
			suffixLen = 2
		} else if strings.HasSuffix(maxSize, "mb") {
			suffixLen = 2
		} else if strings.HasSuffix(maxSize, "kb") {
			unit = 1024
			suffixLen = 2
		} else if strings.HasSuffix(maxSize, "g") {
			unit = 1024 * unit
			suffixLen = 1
		} else if strings.HasSuffix(maxSize, "m") {
			suffixLen = 1
		} else if strings.HasSuffix(maxSize, "k") {
			unit = 1024
			suffixLen = 1
		}
		maxSize = maxSize[:len(maxSize)-suffixLen]
		maxSizeInt, err := strconv.Atoi(maxSize)
		if err == nil {
			pi.MaxSize = int64(maxSizeInt * unit)
		}
	}
	targetFileNamePattern, ok20 := (*cfg)["zip-package-target-file-name-pattern"]
	if ok20 {
		pi.TargetFileNamePattern = targetFileNamePattern
	}
	sourceID, ok30 := (*cfg)["zip-package-source-id"]
	if ok30 {
		pi.SourceID = sourceID
	}
	metaXmlFileName, ok40 := (*cfg)["zip-package-meta-xml-file-name"]
	if ok40 {
		pi.MetaXmlFileName = metaXmlFileName
	}
}

func configParseInstructure(pi *service.ParseInstruction, cfg *map[string]string) {
	//obtain parse instruction info from config to set the following values
	//pi.SetGroupNameDelimiter("[", "]")
	//pi.SetGroupIdNameDelimiter(":")
	//pi.SetContinuousEmptyColLimit(10)
	//pi.SetContinuousEmptyRowLimit(10)
	grpNamePrefix, ok10 := (*cfg)["group-name-prefix"]
	grpNameSuffix, ok20 := (*cfg)["group-name-suffix"]
	if ok10 && ok20 {
		pi.SetGroupNameDelimiter(grpNamePrefix, grpNameSuffix)
	}
	grpIdNameDelim, ok30 := (*cfg)["group-id-name-delimiter"]
	if ok30 {
		pi.SetGroupIdNameDelimiter(grpIdNameDelim)
	}
	emptyColLimit, ok40 := (*cfg)["continuous-empty-col-limit"]
	if ok40 {
		emptyColLimitInt, err := strconv.Atoi(emptyColLimit)
		if err == nil {
			pi.SetContinuousEmptyColLimit(int8(emptyColLimitInt))
		}
	}
	emptyRowLimit, ok50 := (*cfg)["continuous-empty-row-limit"]
	if ok50 {
		emptyRowLimitInt, err := strconv.Atoi(emptyRowLimit)
		if err == nil {
			pi.SetContinuousEmptyRowLimit(int8(emptyRowLimitInt))
		}
	}
}

func reconcile(reportDir, fileEndsWith, outDir, xls, config, sheetName string) {
	fmt.Printf("reconcile: %s %s %s %s %s %s\n", reportDir, outDir, fileEndsWith, xls, config, sheetName)
	cfg := loadConfig(config)
	pi := service.NewParseInstruction()
	pi.SheetName = sheetName
	configParseInstructure(pi, cfg)
	pkg, err := pi.ParsePackageRequests(xls)
	if err != nil || pkg == nil {
		fmt.Printf("ParsePackageExcel failed or no requests at all: %v\n", err)
		return
	}
	fmt.Println("Parse success and get requests: ", len(pkg.Requests))
	ri := service.NewReconcileInstruction()
	ri.ReportDir = reportDir
	ri.OutDir = outDir
	ri.ReportFileEndsWith = fileEndsWith
	reconcileResults, _ := ri.Reconcile(pkg)
	fmt.Println("Reconcile success and get results: ", len(*reconcileResults))
	colHeaders := pi.ExtractRequestHeaders(xls)
	outXls := excelize.NewFile()
	ri.OutputExcel(reconcileResults, colHeaders, outXls)
	targetFile := outDir + "/reconcile-result--" + sheetName + ".xlsx"
	fmt.Println("Output to: ", targetFile)
	outXls.SaveAs(targetFile)
}

func usage() {
	fmt.Printf("Usage: %s --command package --file-dir path/to/input-files --out-dir path/to/output-zip --xls path/to/meta-excel-file --config path/to/config-file --sheet-name default-1st-sheet --unzip true-or-false\n", os.Args[0])
	fmt.Printf("     : %s --command reconcile --report-dir path/to/report --report-file-ends-with .xml --out-dir path/to/reconcile-report --xls path/to/meta-excel-file --config path/to/config-file --sheet-name default-1st-sheet\n", os.Args[0])
}

func loadConfig(cfgFile string) *map[string]string {
	cfg := make(map[string]string)
	f, _ := os.Open(cfgFile)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		cfg[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return &cfg
}
