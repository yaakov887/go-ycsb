package measurement

import (
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

func (r *rawseries) MeasureExt(op string, start time.Time, end time.Time, key string, values []interface{}) {
	rm := rawmeasurement{
		opType:  op,
		opStart: start,
		opEnd:   end,
		opKey:   key,
		opVals:  values,
	}

	r.series = append(r.series, rm)
}
