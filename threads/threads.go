package main

import (
	"github.com/guonaihong/flag"
)

type threads struct {
	maxThreads   int
	maxRequests  int
	threadYields int
	threadlocks  int
}

func (t *threads) run() {
}

func Main(name string, args []string) {
	commandlLine := flag.NewFlagSet(name, flag.ExitOnError)
	maxThreads := commandlLine.Int("num-threads", 1, "number of threads to use")
	maxRequests := commandlLine.Int("max-requests", 10000, "limit for total number of requests")
	threadYields := commandlLine.Int("thread-yields", 1000, "number of yields to do per request")
	threadlocks := commandlLine.Int("thread-locks", 8, "number of locks per thread")

	commandlLine.Author("guonaihong https://github.com/guonaihong/sysbench2")
	commandlLine.Parse(args)
}
