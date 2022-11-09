package workload

import (
	"encoding/json"
	"os"
	"sort"
)

type Action struct {
	NodeID  string `json:"nodeid"`
	Command string `json:"cmd"`
}

type Event struct {
	RelativeTime int      `json:"time"`
	Actions      []Action `json:"actions"`
}

type EventList []Event

type EventWorkload struct {
	Events EventList `json:"events"`
}

// Len Sort interface implementation so we can sort by RelativeTime
func (s EventList) Len() int {
	return len(s)
}

// Less Sort interface implementation so we can sort by RelativeTime
func (s EventList) Less(i, j int) bool {
	return s[i].RelativeTime < s[j].RelativeTime
}

// Swap Sort interface implementation so we can sort by RelativeTime
func (s EventList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// ParseEventList reads the events file to internal data structure
func ParseEventList(jsonSource string) error {
	bytes, err := os.ReadFile(jsonSource)
	if err != nil {
		return err
	}

	var tempList EventWorkload
	err = json.Unmarshal(bytes, &tempList)
	if err != nil {
		return err
	}

	globalEventWorkload = tempList
	sort.Sort(globalEventWorkload.Events)

	return nil
}

var globalEventWorkload EventWorkload
