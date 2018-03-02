package grpool_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/silentred/toolkit/util/grpool"
)

func increment() {
	for i := 0; i < 100; i++ {
	}
}

func Test_GrpoolMemUsage(t *testing.T) {
	n := 100
	for i := 0; i < n; i++ {
		grpool.Add(increment)
	}
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	fmt.Println("mem usage:", mem.TotalAlloc/1024)
}

func Test_GroroutineMemUsage(t *testing.T) {
	n := 100
	for i := 0; i < n; i++ {
		go increment()
	}
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	fmt.Println("mem usage:", mem.TotalAlloc/1024)
}
