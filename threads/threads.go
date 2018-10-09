package threads

import (
	"github.com/guonaihong/flag"
	"runtime"
	"sync"
	"sync/atomic"
)

type threads struct {
	maxThreads   int
	maxRequests  int
	threadYields int
	threadlocks  int
	count        int32
	mutex        []sync.Mutex
}

func (t *threads) run() {

	work := make(chan struct{}, 100)
	var wg sync.WaitGroup

	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < t.maxRequests; i++ {
			work <- struct{}{}
		}
		close(work)
	}()

	wg.Add(t.maxThreads)
	for i := 0; i < t.maxThreads; i++ {
		go func() {
			defer wg.Done()
			for range work {
				newCount := atomic.AddInt32(&t.count, 1)
				newCount = newCount % int32(t.threadlocks)
				for j := 0; i < t.threadYields; j++ {
					t.mutex[newCount].Lock()
					runtime.Gosched()
					t.mutex[newCount].Unlock()
				}
			}
		}()
	}
}

func Main(name string, args []string) {
	commandlLine := flag.NewFlagSet(name, flag.ExitOnError)
	maxThreads := commandlLine.Int("num-threads", 1, "number of threads to use")
	maxRequests := commandlLine.Int("max-requests", 10000, "limit for total number of requests")
	threadYields := commandlLine.Int("thread-yields", 1000, "number of yields to do per request")
	threadlocks := commandlLine.Int("thread-locks", 8, "number of locks per thread")

	commandlLine.Author("guonaihong https://github.com/guonaihong/sysbench2")
	commandlLine.Parse(args)

	t := threads{
		maxThreads:   *maxThreads,
		maxRequests:  *maxRequests,
		threadYields: *threadYields,
		threadlocks:  *threadlocks,
	}

	t.run()
}
