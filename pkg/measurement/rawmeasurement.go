package measurement

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/pingcap/go-ycsb/pkg/ycsb"
	"golang.org/x/crypto/openpgp/errors"
	"strings"
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

func newRawSeries(p *properties.Properties) *rawseries {
	r := new(rawseries)
	return r
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

func (r *rawseries) GetMeasurement(index int) ([]string, error) {
	if len(r.series) == 0 || index > len(r.series) || index < 0 {
		return nil, errors.InvalidArgumentError(index)
	}
	line := []string{}
	line = append(line, r.series[index].opType)
	line = append(line, r.series[index].opStart.String())
	line = append(line, r.series[index].opEnd.String())
	line = append(line, r.series[index].opKey)
	var vals []string
	for v := range r.series[index].opVals {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	line = append(line, strings.Join(vals, ","))

	return line, nil
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
