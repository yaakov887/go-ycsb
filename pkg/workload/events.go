package workload

import (
	"encoding/json"
	"fmt"
	"github.com/pingcap/go-ycsb/pkg/nodectrl"
	"log"
	"os"
	"sort"
	"time"
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

// executeAllActions runs the command on the node specified for all actions in the event's list
func (e *Event) executeAllActions() {
	fmt.Printf("Executing node actions (Count:%v)\n", len(e.Actions))
	for _, a := range e.Actions {
		fmt.Printf("[executeAllActions] (%v:%v)\n", a.NodeID, a.Command)
		err := nodectrl.RunNodeCommand(a.NodeID, a.Command)
		if err != nil {
			log.Printf("ERROR [executeAllActions] (%v:%v) - %v\n", a.NodeID, a.Command, err.Error())
		}
	}
}

// StartEventWorkload spins off multiple go routines to execute the events
// at the time relative to the start of the workload
func StartEventWorkload(jsonSource string) error {
	var err error
	if globalEventWorkload.Events == nil || len(globalEventWorkload.Events) <= 0 {
		err = ParseEventList(jsonSource)
	}
	if err != nil {
		return err
	}

	for _, event := range globalEventWorkload.Events {
		fmt.Printf("Spinning off event {Relative Time:%v, Action Count:%v}\n",
			event.RelativeTime, len(event.Actions))
		ticker := time.NewTicker(time.Duration(event.RelativeTime) * time.Second)
		go func(e Event, t *time.Ticker) {
			for {
				select {
				case execTime := <-t.C:
					fmt.Printf("%v Executing event {Relative Time:%v, Action Count:%v}\n",
						execTime, e.RelativeTime, len(e.Actions))
					e.executeAllActions()
					t.Stop()
					return
				}
			}
		}(event, ticker)
	}

	return nil
}

var globalEventWorkload EventWorkload
