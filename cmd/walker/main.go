package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	walker "github.com/arzonus/adjust_test_task"
)

var (
	flagParallel = flag.Uint("parallel", 10, "control count of parallel requests")
	flagTimeout  = flag.Duration("timeout", 30*time.Second, "timeout of http requests")
)

func main() {
	flag.Parse()

	httpClient := &http.Client{
		Timeout: *flagTimeout,
	}

	if err := walker.NewWalker(*flagParallel, httpClient).Walk(flag.Args()...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
