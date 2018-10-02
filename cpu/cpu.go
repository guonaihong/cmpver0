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
	start int
	end   int
}

type cpu struct {
	task       chan prime
	maxTime    time.Duration
	maxThreads int
	maxPrime   int
}

func calPrime(start, end int) {
	for i := start; i < end; i++ {
		t := int(math.Sqrt(float64(i)))

		l := 2
		for ; i <= t; l++ {
			if t%l == 0 {
				break
			}
		}

		if l > t {
			fmt.Printf("%d\n", t)
		}
	}
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

		step := c.maxPrime % c.maxThreads
		for i := 0; i < c.maxPrime; i += step {

			end := i + step
			if end > c.maxPrime {
				end = c.maxPrime
			}

			c.task <- prime{start: i, end: end}
		}
	}()

	for i := 0; i < c.maxThreads; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range c.task {
				calPrime(v.start, v.end)
			}

		}()
	}
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
}
