package memory

import (
	"fmt"
	"github.com/guonaihong/flag"
	"github.com/guonaihong/gutil/file"
	"strings"
	"sync/atomic"
)

type operType int

const (
	read = iota
	write
	none
)

type modeType int

const (
	seq = iota
	rnd
)

type memory struct {
	buffer     []int64
	blockSize  int
	totalSize  int
	maxThreads int
	memoryOper operType
	accessMode modeType
}

func operS2N(oper string) (operType, error) {
	switch strings.ToLower(oper) {
	case "read":
		return read, nil
	case "write":
		return write, nil
	case "none":
		return none, nil
	default:
		return none, fmt.Errorf("Invalid value for memory-oper: %s", oper)
	}
}

func modeS2N(mode string) (modeType, error) {
	switch strings.ToLower(mode) {
	case "seq":
		return seq, nil
	case "rnd":
		return rnd, nil
	default:
		return rnd, fmt.Errorf("Invalid value for memory-access-mode: %s", mode)
	}
}

func (m *memory) rndNone() {
	for i, l := 0, len(m.buffer); i < l; i++ {
	}
}

func (m *memory) seqNone() {
	for i, l := 0, len(m.buffer); i < l; i++ {
	}
}

func (m *memory) rndRead() {
	for i, l := 0, len(m.buffer); i < l; i++ {
		atomic.LoadInt64(&m.buffer[i])
	}
}

func (m *memory) seqRead() {
}

func (m *memory) rndWrite() {
}

func (m *memory) seqWrite() {
}

func (m *memory) run() {
}

func Main(name string, args []string) {
	commandlLine := flag.NewFlagSet(name, flag.ExitOnError)

	maxThreads := commandlLine.Int("num-threads", 1, "number of threads to use")
	blockSize := commandlLine.String("memory-block-size", "1k", "size of memory block for test")
	totalSize := commandlLine.String("memory-total-size", "100G", "total size of data to transfer")
	memoryOper := commandlLine.String("memory-oper", "write", "type of memory operations {read, write, none}")
	accessMode := commandlLine.String("memory-access-mode", "seq", "memory access mode {seq,rnd}")

	commandlLine.Parse(args)

	mem := memory{
		maxThreads: *maxThreads,
	}

	if oper, err := operS2N(*memoryOper); err != nil {
		fmt.Printf("%s\n", err)
		commandlLine.Usage()
		return
	} else {
		mem.memoryOper = oper
	}

	if mode, err := modeS2N(*accessMode); err != nil {
		fmt.Printf("%s\n", err)
		commandlLine.Usage()
		return
	} else {
		mem.accessMode = mode
	}

	if blockSize0, err := file.ParseSize(*blockSize); err != nil {
		fmt.Printf("Invalid value for memory-block-size %s\n", *blockSize)
		commandlLine.Usage()
		return
	} else {
		mem.blockSize = int(blockSize0)
		mem.buffer = make([]int64, blockSize0/8)
	}

	if totalSize0, err := file.ParseSize(*totalSize); err != nil {
		fmt.Printf("Invalid value for memory-total-size %s\n", *totalSize)
		commandlLine.Usage()
		return
	} else {
		mem.totalSize = int(totalSize0)
	}

	mem.run()
}
