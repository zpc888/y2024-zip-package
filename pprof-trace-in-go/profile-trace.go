package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to file")
var tracefile = flag.String("tracefile", "", "write trace to file")

func main() {
	doWork := func() {
		for i := 0; i < 1000000; i++ {
			n := rand.Intn(100)
			s := n * n * n
			fmt.Printf("Cube of %d is %d\n", n, s)
		}
	}

	flag.Parse()

	// Start CPU profiling
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	// Start tracing
	if *tracefile != "" {
		traceFile, err := os.Create(*tracefile)
		if err != nil {
			panic(err)
		}
		defer traceFile.Close()

		if err := trace.Start(traceFile); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	doWork()

	// Memory profiling
	if *memprofile != "" {
		memF, err := os.Create(*memprofile)
		if err != nil {
			panic(err)
		}
		defer memF.Close()

		if err := pprof.WriteHeapProfile(memF); err != nil {
			panic(err)
		}
		fmt.Println("Memory profile written to mem.prof")

		time.Sleep(5 * time.Second)
	}

	fmt.Println("Done")
}
