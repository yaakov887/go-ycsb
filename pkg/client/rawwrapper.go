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

package client

import (
	"context"
	"fmt"
	"time"

	"github.com/pingcap/go-ycsb/pkg/measurement"
	"github.com/pingcap/go-ycsb/pkg/ycsb"
)

// RawWrapper stores the pointer to a implementation of ycsb.DB.
type RawWrapper struct {
	DB ycsb.DB
}

func rawmeasure(start time.Time, end time.Time, op string, key string, values []interface{}, err error) {
	if err != nil {
		measurement.RawMeasure(fmt.Sprintf("%s_ERROR", op), start, end, key, values)
		return
	}

	measurement.RawMeasure(op, start, end, key, values)
}

func (db RawWrapper) Close() error {
	return db.DB.Close()
}

func (db RawWrapper) InitThread(ctx context.Context, threadID int, threadCount int) context.Context {
	return db.DB.InitThread(ctx, threadID, threadCount)
}

func (db RawWrapper) CleanupThread(ctx context.Context) {
	db.DB.CleanupThread(ctx)
}

func (db RawWrapper) Read(ctx context.Context, table string, key string, fields []string) (_ map[string][]byte, err error) {
	start := time.Now()
	dbRead, err := db.DB.Read(ctx, table, key, fields)
	end := time.Now()
	var tempVals []interface{}
	for _, dbVal := range dbRead {
		tempVals = append(tempVals, dbVal)
	}
	rawmeasure(start, end, "READ", key, tempVals, err)

	return dbRead, err
}

func (db RawWrapper) BatchRead(ctx context.Context, table string, keys []string, fields []string) (_ []map[string][]byte, err error) {
	return nil, nil
}

func (db RawWrapper) Scan(ctx context.Context, table string, startKey string, count int, fields []string) (_ []map[string][]byte, err error) {
	return nil, nil
}

func (db RawWrapper) Update(ctx context.Context, table string, key string, values map[string][]byte) (err error) {
	var tempVals []interface{}
	for _, pVal := range values {
		tempVals = append(tempVals, pVal)
	}

	start := time.Now()
	err = db.DB.Update(ctx, table, key, values)
	rawmeasure(start, time.Now(), "UPDATE", key, tempVals, err)
	return err
}

func (db RawWrapper) BatchUpdate(ctx context.Context, table string, keys []string, values []map[string][]byte) (err error) {
	return nil
}

func (db RawWrapper) Insert(ctx context.Context, table string, key string, values map[string][]byte) (err error) {
	start := time.Now()
	var tempVals []interface{}
	for _, pVal := range values {
		tempVals = append(tempVals, pVal)
	}

	defer func() {
		rawmeasure(start, time.Now(), "INSERT", key, tempVals, err)
	}()

	return db.DB.Insert(ctx, table, key, values)
}

func (db RawWrapper) BatchInsert(ctx context.Context, table string, keys []string, values []map[string][]byte) (err error) {
	return nil
}

func (db RawWrapper) Delete(ctx context.Context, table string, key string) (err error) {
	return nil
}

func (db RawWrapper) BatchDelete(ctx context.Context, table string, keys []string) (err error) {
	return nil
}

func (db RawWrapper) Analyze(ctx context.Context, table string) error {
	return nil
}
