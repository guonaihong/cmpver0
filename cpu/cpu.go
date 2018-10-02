package cpu

import (
	"fmt"
	"github.com/guonaihong/flag"
	time2 "github.com/guonaihong/gutil/time"
	"math"
	"sync"
	"time"
)

type prime struct {
	i int
}

type cpu struct {
	task       chan prime
	maxTime    time.Duration
	maxThreads int
	maxPrime   int
}

func calPrime(i int) {
	if i < 3 {
		i = 3
	}

	t := int(math.Sqrt(float64(i)))

	j := 2
	for ; j <= t; j++ {
		if i%j == 0 {
			return
		}
	}

	//fmt.Printf("%d\n", i)
}

func (c *cpu) run() {
	var wg sync.WaitGroup

	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer func() {
			close(c.task)
			wg.Done()
		}()

		for i := 3; i < c.maxPrime; i++ {

			c.task <- prime{i: i}
		}
	}()

	for i := 0; i < c.maxThreads; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range c.task {
				calPrime(v.i)
			}

		}()
	}
}

func (c *cpu) report() {
	fmt.Printf("Doing CPU performance benchmark\n")
}

func Main(name string, args []string) {
	commandlLine := flag.NewFlagSet(name, flag.ExitOnError)
	maxPrime := commandlLine.Int("cpu-max-prime", 10000, "upper limit for primes generator")
	maxTime := commandlLine.String("max-time", "0", "limit for total execution time in seconds")
	maxThreads := commandlLine.Int("num-threads", 1, "number of threads to use")
	commandlLine.Author("guonaihong https://github.com/guonaihong/sysbench2")
	commandlLine.Parse(args)

	if *maxPrime < 0 {
		fmt.Printf("Invalid value of cpu-max-prime: %d.\n", *maxPrime)
		return
	}

	c := cpu{
		task:       make(chan prime, 128),
		maxPrime:   *maxPrime,
		maxTime:    time2.ParseTime(*maxTime),
		maxThreads: *maxThreads,
	}

	c.run()
	c.report()
}
