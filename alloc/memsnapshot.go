package alloc

import (
	"runtime"

	"github.com/shirou/gopsutil/mem"
)

type MemSnapshot struct {
	Used         int64
	Sys          int64
	HeapSys      int64
	HeapIdle     int64
	HeapInuse    int64
	HeapReleased int64
}

func NewMemSnapshot() MemSnapshot {
	vm, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return MemSnapshot{
		Used:         int64(vm.Used),
		Sys:          int64(memStats.Sys),
		HeapSys:      int64(memStats.HeapSys),
		HeapIdle:     int64(memStats.HeapIdle),
		HeapInuse:    int64(memStats.HeapInuse),
		HeapReleased: int64(memStats.HeapReleased),
	}
}
