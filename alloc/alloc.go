package alloc

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/mem"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Alloc struct {
	cfg Config

	amount int64
	count  int64

	csvWriter *csv.Writer
}

func NewAlloc(cfg Config) *Alloc {
	alloc := Alloc{
		cfg: cfg,
	}
	alloc.csvWriter = alloc.initCsv(fmt.Sprintf("M%dmin%dmax%ds%vrmin%drmax%drs%v.csv",
		cfg.MaxLimit, cfg.MaxSize, cfg.MinSize, cfg.Spread, cfg.ReMinSize, cfg.ReMaxSize, cfg.ReSpread))
	return &alloc
}

func (a *Alloc) Run() {
	bufs := make([][]byte, 0)

	for {
		size := a.size(a.cfg.MinSize, a.cfg.MaxSize, a.cfg.Spread)
		bufs = append(bufs, a.alloc(size))
		a.count = int64(len(bufs))
		a.amount += size

		a.writeLine()

		if a.amount >= a.cfg.MaxLimit {
			break
		}
	}

	bufs = bufs[:0]
	a.amount = 0
	a.count = 0
	runtime.GC()

	for {
		size := a.size(a.cfg.ReMinSize, a.cfg.ReMaxSize, a.cfg.ReSpread)
		bufs = append(bufs, a.alloc(size))
		a.count = int64(len(bufs))
		a.amount += size

		a.writeLine()

		if a.amount >= 2*a.cfg.MaxLimit {
			break
		}
	}
}

func (a *Alloc) size(min, max int64, spread float64) int64 {
	r := rand.Float64()

	m := int64(float64(min+max) * (1 - spread))

	if r < spread {
		return rand.Int63()%(m-min) + min
	} else {
		return rand.Int63()%(max-m) + m
	}
}

func (a *Alloc) alloc(size int64) []byte {
	buf := make([]byte, size)
	for i := 0; i < len(buf); i += 4095 {
		buf[i] = byte(i % 127)
	}
	return buf
}

func (a *Alloc) initCsv(filename string) *csv.Writer {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	csvWriter := csv.NewWriter(file)

	header := []string{
		"amount",
		"count",

		"total",
		"used",
		"usedPercent",
		"cached",
		"slab",

		"alloc",
		"totalAlloc",
		"sys",
		"heapSys",
		"mallocs",
		"frees",
		"heapAlloc",
		"heapIdle",
		"heapInuse",
		"heapReleased",
		"heapObjects",
		"stackInuse",
		"stackSys",
		"mspanInuse",
		"mspanSys",
		"mcacheInuse",
		"mcacheSys",
		"ohterSys",
	}
	if err := csvWriter.Write(header); err != nil {
		panic(err)
	}
	return csvWriter
}

func (a *Alloc) writeLine() {
	vm, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	str := []string{
		formatMemused(uint64(a.amount)),
		fmt.Sprintf("%v", a.count),
		formatMemused(vm.Total),
		formatMemused(vm.Used),
		fmt.Sprintf("%v", vm.UsedPercent),
		formatMemused(vm.Cached),
		formatMemused(vm.Slab),
		formatMemused(memStats.Alloc),
		formatMemused(memStats.TotalAlloc),
		formatMemused(memStats.Sys),
		formatMemused(memStats.HeapSys),
		fmt.Sprintf("%v", memStats.Mallocs),
		fmt.Sprintf("%v", memStats.Frees),
		formatMemused(memStats.HeapAlloc),
		formatMemused(memStats.HeapIdle),
		formatMemused(memStats.HeapInuse),
		formatMemused(memStats.HeapReleased),
		fmt.Sprintf("%v", memStats.HeapObjects),
		formatMemused(memStats.StackInuse),
		formatMemused(memStats.StackSys),
		formatMemused(memStats.MSpanInuse),
		formatMemused(memStats.MSpanSys),
		formatMemused(memStats.MCacheInuse),
		formatMemused(memStats.MCacheSys),
		formatMemused(memStats.OtherSys),
	}

	if err := a.csvWriter.Write(str); err != nil {
		panic(err)
	}
	a.csvWriter.Flush()

}

func formatMemused(b uint64) string {
	return fmt.Sprintf("%.2f", float64(b)/float64(1024*1024))
}
