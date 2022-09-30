package measurement

import (
	"github.com/pingcap/go-ycsb/pkg/ycsb"
	"time"
)

type rawmeasurement struct {
	opType  string
	opStart time.Time
	opEnd   time.Time
	opKey   string
	opVals  []interface{}
}

type rawseries struct {
	series []rawmeasurement
}

func (r *rawseries) Measure(op string, start time.Time, end time.Time, key string, values []interface{}) {
	rm := rawmeasurement{
		opType:  op,
		opStart: start,
		opEnd:   end,
		opKey:   key,
		opVals:  values,
	}

	r.series = append(r.series, rm)
}

// Summary returns the summary of the measurement.
func (r *rawseries) Summary() []string {
	return nil
}

func (r *rawseries) Info() ycsb.MeasurementInfo {
	tempInfo := make(map[string]interface{})
	tempInfo["len"] = len(r.series)
	return newRawmeasurementInfo(tempInfo)
}

type rawmeasurementInfo struct {
	info map[string]interface{}
}

func newRawmeasurementInfo(info map[string]interface{}) *rawmeasurementInfo {
	return &rawmeasurementInfo{info: info}
}

func (ri *rawmeasurementInfo) Get(metricName string) interface{} {
	if value, ok := ri.info[metricName]; ok {
		return value
	}
	return nil
}
