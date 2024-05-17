package main

import (
	"encoding/xml"
	"fmt"
	"zip-pkg-in-go/model"
)

func main() {
	fmt.Println("Hello, World!")
	ingestSource := &model.Attribute{
		Name:  "ingestedSource",
		Value: "s3://bucket-name/folder-name/",
	}
	out, _ := xml.MarshalIndent(ingestSource, "", "  ")
	fmt.Println(string(out))

}
