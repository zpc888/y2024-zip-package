package main

import (
	_ "flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Printf("Usage: %s --command package --file-dir path/to/input-files --xls path/to/meta-excel-file --config path/to/config-file --sheet-index default-1st-sheet\n", os.Args[0])
		fmt.Printf("     : %s --command reconcile --report-dir path/to/report --xls path/to/meta-excel-file --config path/to/config-file --sheet-index default-1st-sheet\n", os.Args[0])
		os.Exit(1)
	}
}
