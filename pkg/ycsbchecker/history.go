package ycsbchecker

// This code was initially implemented in the ailidani-paxi go module
// Modified for the Charapko YCSB Event Based Performance Benchmark

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
)

// History client operation history mapped by key
type History struct {
	sync.RWMutex
	shard      map[string][]*operation
	operations []*operation
}

// NewHistory creates a History map
func NewHistory() *History {
	return &History{
		shard:      make(map[string][]*operation),
		operations: make([]*operation, 0),
	}
}

// Add puts an operation in History
func (h *History) Add(key string, input, output interface{}, start, end int64) {
	h.Lock()
	defer h.Unlock()
	if _, exists := h.shard[key]; !exists {
		h.shard[key] = make([]*operation, 0)
	}
	o := &operation{input, output, start, end}
	h.shard[key] = append(h.shard[key], o)
	h.operations = append(h.operations, o)
}

// AddOperation adds the operation
func (h *History) AddOperation(key string, o *operation) {
	h.Lock()
	defer h.Unlock()
	if _, exists := h.shard[key]; !exists {
		h.shard[key] = make([]*operation, 0)
	}
	h.shard[key] = append(h.shard[key], o)
	h.operations = append(h.operations, o)
}

// Linearizable concurrently checks if each partition of the history is linearizable and returns the total number of anomaly reads
func (h *History) Linearizable() int {
	anomalies := make(chan []*operation)
	h.RLock()
	defer h.RUnlock()
	for _, partition := range h.shard {
		c := newChecker()
		go func(p []*operation) {
			anomalies <- c.linearizable(p)
		}(partition)
	}
	sum := 0
	for range h.shard {
		a := <-anomalies
		sum += len(a)
	}
	return sum
}

// WriteFile writes entire operation history into file
func (h *History) WriteFile(path string) error {
	file, err := os.Create(path + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	h.RLock()
	defer h.RUnlock()

	sort.Sort(byTime(h.operations))

	latency := 0.0
	throughput := 0
	s := 1.0
	for _, o := range h.operations {
		start := float64(o.start) / 1000000000.0
		end := float64(o.end) / 1000000000.0
		fmt.Fprintf(w, "%v,%v,%f,%f\n", o.input, o.output, start, end)
		latency += end - start
		throughput++
		if end > s {
			fmt.Fprintf(w, "PerSecond %f %d\n", latency/float64(throughput)*1000.0, throughput)
			latency = 0
			throughput = 0
			s++
		}

		// fmt.Fprintln(w, o)
	}

	// for k, ops := range h.shard {
	// 	fmt.Fprintf(w, "key=%d\n", k)
	// 	for _, o := range ops {
	// 		fmt.Fprintln(w, o)
	// 	}
	// }
	return w.Flush()
}

// ReadFile reads csv log file and create operations in history
func (h *History) ReadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(record) < 5 {
			return errors.New("operation history file format error")
		}

		// get id / key
		//{"Operation", "Start", "End", "Key", "Value(s)"}
		id := record[3]

		operation := new(operation)

		switch record[0] {
		case "Operation":
			continue
		case "READ":
			operation.input = nil
			operation.output = record[4]
		case "UPDATE":
		case "INSERT":
			operation.input = record[4]
			operation.output = nil
		default:
			operation.input = nil
			operation.output = nil
		}

		// get start time
		start, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Fatal(err)
			return err
		}
		operation.start = start

		// get end time
		end, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			log.Fatal(err)
			return err
		}
		operation.end = end

		h.AddOperation(id, operation)
	}

	return file.Close()
}
