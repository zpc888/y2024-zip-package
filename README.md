# Overview
This tool is to generate a list of zip file based on the configured limitation for each zip and a list of source files with their own metadata in an excel spreadsheet. It will provide java and go implementation.
# Java implementation
Provide a native image via GraalVM to avoid JDK installation on client box.
# Go implementation
Try to see whether Go is better than Java performance-wise with normal development thought.

After `go build`, the tool size is 6.1M in Ubuntu box vs 60M in java GraalVM native image.
