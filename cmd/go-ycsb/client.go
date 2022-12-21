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

package main

import (
	"fmt"
	"github.com/pingcap/go-ycsb/pkg/client"
	"github.com/pingcap/go-ycsb/pkg/measurement"
	"github.com/pingcap/go-ycsb/pkg/nodectrl"
	"github.com/pingcap/go-ycsb/pkg/workload"
	"github.com/pingcap/go-ycsb/pkg/ycsbchecker"
	"strconv"
	"time"

	"github.com/pingcap/go-ycsb/pkg/prop"
	"github.com/spf13/cobra"
)

func runNodeCommandFunc(cmd *cobra.Command, args []string, wt string) {
	switch wt {
	case "startnodes":
		runStartNodesCommandFunc(cmd, args)
	case "stopnodes":
		runStopNodesCommandFunc(cmd, args)
	}
}

func runCoreWorkloadCommandFunc() {
	fmt.Println("***************** properties *****************")
	for key, value := range globalProps.Map() {
		fmt.Printf("\"%s\"=\"%s\"\n", key, value)
	}
	fmt.Println("**********************************************")

	var err error
	if !nodectrl.NodesParsed() {
		nodeSrc := globalProps.GetString(prop.Cluster, "")
		err = nodectrl.ParseNodeList(nodeSrc)
	}
	if err == nil {
		eventSrc := globalProps.GetString(prop.Events, "")
		if eventSrc != "" {
			err = workload.StartEventWorkload(eventSrc)
			if err != nil {
				fmt.Printf("Error creating workload events [%v]\n", err.Error())
			}
		}
	} else {
		fmt.Printf("Error parsing node info [%v]\n", err.Error())
	}

	c := client.NewClient(globalProps, globalWorkload, globalDB)
	start := time.Now()
	c.Run(globalContext)

	fmt.Printf("Run finished, takes %s\n", time.Now().Sub(start))
	measurementType, _ := globalProps.Get(prop.MeasurementType)
	if measurementType == "raw" {
		measurement.RawOutput()
	} else {
		measurement.Output()
	}

	checkerType := globalProps.GetString(prop.Checker, "")
	filePrefix := globalProps.GetString(prop.CSVFileName, "")
	if checkerType != "" {
		err = ycsbchecker.RunChecker(checkerType, filePrefix)
	}
}

func runClientCommandFunc(cmd *cobra.Command, args []string, doTransactions bool, command string) {
	time.Sleep(30 * time.Second)
	dbName := args[0]

	initialGlobal(dbName, func() {
		doTransFlag := "true"
		if !doTransactions {
			doTransFlag = "false"
		}
		globalProps.Set(prop.DoTransactions, doTransFlag)
		globalProps.Set(prop.Command, command)

		if cmd.Flags().Changed("threads") {
			// We set the threadArg via command line.
			globalProps.Set(prop.ThreadCount, strconv.Itoa(threadsArg))
		}

		if cmd.Flags().Changed("target") {
			globalProps.Set(prop.Target, strconv.Itoa(targetArg))
		}

		if cmd.Flags().Changed("interval") {
			globalProps.Set(prop.LogInterval, strconv.Itoa(reportInterval))
		}
	})

	workType := globalProps.GetString(prop.Workload, "core")

	switch workType {
	case "core":
		runCoreWorkloadCommandFunc()
	case "startnodes":
	case "stopnodes":
		runNodeCommandFunc(cmd, args, workType)
	default:
		fmt.Printf("[ERROR] Invalid workload property value: %v\n", workType)
	}
}

func runLoadCommandFunc(cmd *cobra.Command, args []string) {
	if len(propertyFiles) > 0 {
		for i, work := range propertyFiles {
			currentWork = work
			fmt.Printf("Loading workload %v %v...\n", i+1, currentWork)
			runClientCommandFunc(cmd, args, false, "load")
			fmt.Printf("Completed workload %v %v...\n", i+1, currentWork)
		}
	} else {
		currentWork = ""
		runClientCommandFunc(cmd, args, false, "load")
	}
}

func runTransCommandFunc(cmd *cobra.Command, args []string) {
	time.Sleep(30 * time.Second)
	if len(propertyFiles) > 0 {
		for i, work := range propertyFiles {
			currentWork = work
			fmt.Printf("Running workload %v %v...\n", i+1, currentWork)
			runClientCommandFunc(cmd, args, true, "run")
			fmt.Printf("Completed workload %v %v...\n", i+1, currentWork)
		}
	} else {
		currentWork = ""
		runClientCommandFunc(cmd, args, true, "run")
	}
}

func runStartNodesCommandFunc(cmd *cobra.Command, args []string) {
	initialGlobalProps(func() {})

	fmt.Println("***************** properties *****************")
	for key, value := range globalProps.Map() {
		fmt.Printf("\"%s\"=\"%s\"\n", key, value)
	}
	fmt.Println("**********************************************")

	nodectrl.ParseNodeList(globalProps.GetString(prop.Cluster, "./cluster.json"))
	nodectrl.StartNodes()
	return
}

func runStopNodesCommandFunc(cmd *cobra.Command, args []string) {
	initialGlobalProps(func() {})

	fmt.Println("***************** properties *****************")
	for key, value := range globalProps.Map() {
		fmt.Printf("\"%s\"=\"%s\"\n", key, value)
	}
	fmt.Println("**********************************************")

	if !nodectrl.NodesParsed() {
		nodectrl.ParseNodeList(globalProps.GetString(prop.Cluster, "./cluster.json"))
	}
	nodectrl.StopNodes()
	return
}

var (
	threadsArg     int
	targetArg      int
	reportInterval int
)

func initClientCommand(m *cobra.Command) {
	m.Flags().StringSliceVarP(&propertyFiles, "property_file", "P", nil, "Specify a property file")
	m.Flags().StringArrayVarP(&propertyValues, "prop", "p", nil, "Specify a property value with name=value")
	m.Flags().StringVar(&tableName, "table", "", "Use the table name instead of the default \""+prop.TableNameDefault+"\"")
	m.Flags().IntVar(&threadsArg, "threads", 1, "Execute using n threads - can also be specified as the \"threadcount\" property")
	m.Flags().IntVar(&targetArg, "target", 0, "Attempt to do n operations per second (default: unlimited) - can also be specified as the \"target\" property")
	m.Flags().IntVar(&reportInterval, "interval", 10, "Interval of outputting measurements in seconds")
}

func initNodeCommand(m *cobra.Command) {
	m.Flags().StringSliceVarP(&propertyFiles, "property_file", "P", nil, "Specify a property file")
	m.Flags().StringArrayVarP(&propertyValues, "prop", "p", nil, "Specify a property value with name=value")
}

func newLoadCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "load db",
		Short: "YCSB load benchmark",
		Args:  cobra.MinimumNArgs(1),
		Run:   runLoadCommandFunc,
	}

	initClientCommand(m)
	return m
}

func newRunCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "run db",
		Short: "YCSB run benchmark",
		Args:  cobra.MinimumNArgs(1),
		Run:   runTransCommandFunc,
	}

	initClientCommand(m)
	return m
}

func newStartNodesCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "startnodes",
		Short: "YCSB start nodes",
		Args:  cobra.NoArgs,
		Run:   runStartNodesCommandFunc,
	}

	initNodeCommand(m)
	return m
}

func newStopNodesCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "stopnodes",
		Short: "YCSB stop nodes",
		Args:  cobra.NoArgs,
		Run:   runStopNodesCommandFunc,
	}

	initNodeCommand(m)
	return m
}
