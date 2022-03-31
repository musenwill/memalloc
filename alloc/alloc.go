package alloc

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
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

	bufs = make([][]byte, 0)
	a.amount = 0
	a.count = 0
	runtime.GC()

	for {
		size := a.size(a.cfg.ReMinSize, a.cfg.ReMaxSize, a.cfg.ReSpread)
		bufs = append(bufs, a.alloc(size))
		a.count = int64(len(bufs))
		a.amount += size

		a.writeLine()

		if a.amount >= a.cfg.MaxLimit {
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
		formatMemSize(uint64(a.amount)),
		strconv.FormatInt(a.count, 10),
		formatMemSize(uint64(vm.Total)),
		formatMemSize(uint64(vm.Used)),
		strconv.FormatFloat(vm.UsedPercent, 'f', 1, 64),

		formatMemSize(vm.Cached),
		formatMemSize(vm.Slab),
		formatMemSize(memStats.Alloc),
		formatMemSize(memStats.TotalAlloc),
		formatMemSize(memStats.Sys),
		formatMemSize(memStats.HeapSys),

		strconv.FormatInt(int64(memStats.Mallocs), 10),
		strconv.FormatInt(int64(memStats.Frees), 10),

		formatMemSize(memStats.HeapAlloc),
		formatMemSize(memStats.HeapIdle),
		formatMemSize(memStats.HeapInuse),
		formatMemSize(memStats.HeapReleased),
		strconv.FormatInt(int64(memStats.HeapObjects), 10),

		formatMemSize(memStats.StackInuse),
		formatMemSize(memStats.StackSys),
		formatMemSize(memStats.MSpanInuse),
		formatMemSize(memStats.MSpanSys),
		formatMemSize(memStats.MCacheInuse),
		formatMemSize(memStats.MCacheSys),
		formatMemSize(memStats.OtherSys),
	}

	if err := a.csvWriter.Write(str); err != nil {
		panic(err)
	}
	a.csvWriter.Flush()
}

func formatMemSize(size uint64) string {
	return strconv.FormatFloat(float64(size)/float64(1024*1024), 'f', 2, 64)
}
