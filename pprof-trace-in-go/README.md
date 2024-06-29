# Debugging Go Code 
## [Debug GoLang with Delve](https://golang.cafe/blog/golang-debugging-with-delve.html)
## [go-delve](https://github.com/go-delve/delve)
## [Debugging Go in Intellij IDEA](https://www.jetbrains.com/help/go/debugging-code.html)
## Using pprof and trace to Diagnose and Fix Performance Issues
### [Profiling Go Programs](https://blog.golang.org/profiling-go-programs)
### [InfoQ](https://www.infoq.com/articles/debugging-go-programs-pprof-trace/)
```shell
go run profile-trace.go -cpuprofile cpu.prof -memprofile mem.prof -tracefile trace.out

go tool pprof -http=:8080 cpu.prof 
go tool pprof -http=:8081 mem.prof 
go tool trace trace.out 
```