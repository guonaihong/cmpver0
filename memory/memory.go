package memory

import "fmt"

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
		return nono, nil
	default:
		return fmt.Errorf("Invalid value for memory-oper: %s", oper)
	}
}

func modeS2N(mode string) (modeType, error) {
	switch strings.ToLower(mode) {
	case "seq":
		return seq, nil
	case "rnd":
		return rnd, nil
	default:
		return fmt.Errorf("Invalid value for memory-access-mode: %s", oper)
	}
}

func Main(name string, args []string) {
	commandlLine := flag.NewFlagSet(name, flag.ExitOnError)

	maxThreads := commandlLine.Int("num-threads", 1, "number of threads to use")
	blockSize := commandlLine.String("memory-block-size", "1k", "size of memory block for test")
	totalSize := commandlLine.String("memory-total-size", "100G", "total size of data to transfer")
	memoryOper := commandlLine.String("memory-oper", "write", "type of memory operations {read, write, none}")
	accessMode := commandlLine.String("memory-access-mode", "seq", "memory access mode {seq,rnd}")

	commandlLine.Parse(args)

	var oper operType
	var mode modeType
	var err error

	if oper, err = operS2N(*memoryOper); err != nil {
		fmt.Printf("%s\n", err)
		commandlLine.Usage()
		return
	}

	if mode, err = modeS2N(*accessMode); err != nil {
		fmt.Printf("%s\n", err)
		commandlLine.Usage()
		return
	}

	mem := memory{
		maxThreads: *maxThreads,
		memoryOper: oper,
		accessMode: mode,
	}

}
