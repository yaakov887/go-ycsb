package ycsbchecker

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// RunChecker calls the designated checker
func RunChecker(checkType, prefix string) error {
	var err error
	switch checkType {
	case "linearizable":
		err = runLinearizable(prefix)
	default:
		return nil
	}
	return err
}

// runLinearizable runs the ailidani-paxi Linearizable checker
func runLinearizable(prefix string) error {
	history := NewHistory()

	f, err := os.Open(".")
	if err != nil {
		return err
	}

	fileList, err := f.Readdir(0)
	if err != nil {
		return err
	}

	fileErrors := 0
	for _, currentFile := range fileList {
		fname := currentFile.Name()
		if strings.Contains(fname, prefix) {
			err = history.ReadFile(fname)
			if err != nil {
				fileErrors += 1
				log.Printf("[LINEARIZABLE] Error reading file %v {%v}", fname, err.Error())
			}
		}
	}

	if fileErrors > 0 {
		err = errors.New(fmt.Sprintf("[ERROR] Linearizable check returned errors for %v files\n", fileErrors))
	} else {
		anomalies := history.Linearizable()
		fmt.Printf("Linearizable check returned %v anomalies\n", anomalies)
		err = nil
	}

	return err
}
