# Overview
This tool is to generate a list of zip file based on the configured limitation for each zip and a list of source files with their own metadata in an excel spreadsheet. It will provide java and go implementation.
# Java implementation
Provide a native image via GraalVM to avoid JDK installation on client box.
# Go implementation
Try to see whether Go is better than Java performance-wise with normal development thought.

After `go build`, the tool size is 6.1M in Ubuntu box vs 60M in java GraalVM native image.

# Performance comparison:
To zip 7000+ PDF files (totally ~8G), if max split zip file is 1G + ~15 attributes for each PDF, here are the results:
- java execution (JDK 17) - about 17 minutes
- native image java execution - about 17 minutes
- Go - about 4 minutes

Java only uses the main thread (single thread), Go doesn't use co-route either.

