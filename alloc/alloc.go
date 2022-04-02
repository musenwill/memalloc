package alloc

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
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

	csvWriter *os.File
}

func NewAlloc(cfg Config) *Alloc {
	alloc := Alloc{
		cfg: cfg,
	}
	if cfg.Print {
		alloc.csvWriter = alloc.initCsv(fmt.Sprintf("M%dmin%dmax%ds%vrmin%drmax%drs%vhalf%v.csv",
			cfg.MaxLimit, cfg.MinSize, cfg.MaxSize, cfg.Spread, cfg.ReMinSize, cfg.ReMaxSize, cfg.ReSpread, cfg.ReleaseHalf))
	}
	return &alloc
}

func (a *Alloc) Run() {
	startSnapshot := NewMemSnapshot()

	var totalAmount int64 = a.cfg.MaxLimit
	linkA := Link{}
	linkB := Link{}
	var lastAmount int64 = 0

	for {
		size := a.size(a.cfg.MinSize, a.cfg.MaxSize, a.cfg.Spread)

		if a.count&0x01 == 0x01 {
			linkA.Push(a.alloc(size))
		} else {
			linkB.Push(a.alloc(size))
		}

		a.count++
		a.amount += size
		if a.amount-lastAmount >= 16*1024 {
			a.writeLine()
			lastAmount = a.amount
		}

		if a.amount >= totalAmount {
			break
		}
	}

	linkA = Link{}
	if !a.cfg.ReleaseHalf {
		linkB = Link{}
		a.count = 0
		a.amount = 0
	} else {
		totalAmount /= 2
		a.count /= 2
		a.amount /= 2
	}
	lastAmount = 0
	runtime.GC()

	halftimeSnapshot := NewMemSnapshot()

	for {
		size := a.size(a.cfg.ReMinSize, a.cfg.ReMaxSize, a.cfg.ReSpread)

		if a.count&0x01 == 0x01 {
			linkA.Push(a.alloc(size))
		} else {
			linkB.Push(a.alloc(size))
		}

		a.count++
		a.amount += size

		if a.amount-lastAmount >= 16*1024 {
			a.writeLine()
			lastAmount = a.amount
		}

		if a.amount >= totalAmount {
			break
		}
	}

	endSnapshot := NewMemSnapshot()

	xUsed := halftimeSnapshot.Used - startSnapshot.Used
	yUsed := endSnapshot.Used - startSnapshot.Used

	xSys := halftimeSnapshot.Sys - startSnapshot.Sys
	ySys := endSnapshot.Sys - startSnapshot.Sys

	xHeapSys := halftimeSnapshot.HeapSys - startSnapshot.HeapSys
	yHeapSys := endSnapshot.HeapSys - startSnapshot.HeapSys

	xRecycle := halftimeSnapshot.HeapIdle - halftimeSnapshot.HeapReleased
	yRecycle := endSnapshot.HeapIdle - endSnapshot.HeapReleased

	var m int64 = 1024 * 1024
	fmt.Printf("    used: %5v		    used: %5v		rate: %.2f\n", xUsed/m, yUsed/m, float64(yUsed)/float64(xUsed))
	fmt.Printf("     sys: %5v		     sys: %5v		rate: %.2f\n", xSys/m, ySys/m, float64(ySys)/float64(xSys))
	fmt.Printf("heap sys: %5v		heap sys: %5v		rate: %.2f\n", xHeapSys/m, yHeapSys/m, float64(yHeapSys)/float64(xHeapSys))
	fmt.Printf(" recycle: %5v		 recycle: %5v		rate: %.2f\n", xRecycle/m, yRecycle/m, float64(yRecycle)/float64(xRecycle))

	if a.cfg.Print {
		a.csvWriter.Sync()
	}
}

func (a *Alloc) size(min, max int64, spread float64) int64 {
	m := int64(float64(min+max) * (1 - spread))

	if m == min || m == max {
		return m
	}

	r := rand.Float64()
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

func (a *Alloc) initCsv(filename string) *os.File {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
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
		"gcSys",
		"numGC",
	}
	if _, err := file.WriteString(strings.Join(header, ",") + "\n"); err != nil {
		panic(err)
	}
	return file
}

func (a *Alloc) writeLine() {
	if !a.cfg.Print {
		return
	}

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
		formatMemSize(memStats.GCSys),
		strconv.FormatInt(int64(memStats.NumGC), 10),
	}

	if _, err := a.csvWriter.WriteString(strings.Join(str, ",") + "\n"); err != nil {
		panic(err)
	}
}

func formatMemSize(size uint64) string {
	return strconv.FormatFloat(float64(size)/float64(1024*1024), 'f', 2, 64)
}
