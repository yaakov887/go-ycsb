// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package measurement

import (
	"fmt"
	"github.com/pingcap/go-ycsb/pkg/util"
	"os"
	"sync"
	"time"

	"github.com/magiconair/properties"
	"github.com/pingcap/go-ycsb/pkg/prop"
	"github.com/pingcap/go-ycsb/pkg/ycsb"
)

var seriesheader = []string{"Operation", "Start", "End", "Key", "Value(s)"}

type series struct {
	sync.RWMutex

	p *properties.Properties

	rawSeries *rawseries
}

func (s *series) measure(op string, start time.Time, end time.Time, key string, values []interface{}) {
	s.RLock()
	ok := s.rawSeries != nil
	s.RUnlock()

	if !ok {
		s.Lock()
		s.rawSeries = newRawSeries()
		s.Unlock()
	}

	(s.rawSeries).Measure(op, start, end, key, values)
}

func (s *series) output() {
	s.RLock()
	defer s.RUnlock()
	//fmt.Printf("%+v\n", s.rawSeries)

	lines := [][]string{}
	var length int
	length = (s.rawSeries).Info().Get("len").(int)
	fmt.Printf("Series Length: %v\n", length)

	for i := 0; i < length; i++ {
		meas, _ := (s.rawSeries).GetMeasurement(i)
		lines = append(lines, meas)
	}
	//fmt.Printf("%+v\n", lines)

	outputStyle := s.p.GetString(prop.OutputStyle, util.OutputStyleCSV)
	filename := s.p.GetString(prop.CSVFileName, prop.Workload)
	fileHandle, err := os.Create(fmt.Sprintf(
		"%v_%v_%v.csv",
		filename,
		time.Now().UnixMilli(),
		s.rawSeries.Info().Get("len"),
	))
	if err != nil {
		fileHandle = nil
	}

	switch outputStyle {
	case util.OutputStylePlain:
		util.RenderString("%-6s - %s\n", seriesheader, lines)
	case util.OutputStyleJson:
		util.RenderJson(seriesheader, lines)
	case util.OutputStyleTable:
		util.RenderTable(seriesheader, lines)
	case util.OutputStyleCSV:
		util.RenderCSV(seriesheader, lines, fileHandle)
	default:
		panic("unsupported outputstyle: " + outputStyle)
	}
	fmt.Printf("Completed output of %v\n", fileHandle)
}

func (s *series) info() ycsb.MeasurementInfo {
	s.RLock()
	defer s.RUnlock()

	return (s.rawSeries).Info()
}

// RawInitMeasure initializes the global measurement.
func RawInitMeasure(p *properties.Properties) {
	globalRawMeasure = new(series)
	globalRawMeasure.p = p
	globalRawMeasure.rawSeries = newRawSeries()
	EnableWarmUp(p.GetInt64(prop.WarmUpTime, 0) > 0)
}

// RawOutput prints the measurement summary.
func RawOutput() {
	var outputReference *rawseries

	globalRawMeasure.Lock()
	outputReference = globalRawMeasure.rawSeries
	globalRawMeasure.rawSeries = newRawSeries()
	globalRawMeasure.Unlock()

	var outputSeries = &series{
		RWMutex:   sync.RWMutex{},
		p:         globalRawMeasure.p,
		rawSeries: outputReference,
	}
	outputSeries.output()
}

// RawMeasure measures the operation.
func RawMeasure(op string, start time.Time, end time.Time, key string, values []interface{}) {
	if IsWarmUpFinished() {
		globalRawMeasure.measure(op, start, end, key, values)
	}
}

// RawInfo returns all the operations MeasurementInfo.
// The key of returned map is the operation name.
func RawInfo() ycsb.MeasurementInfo {
	return globalRawMeasure.info()
}

var globalRawMeasure *series
