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

// DbWrapper stores the pointer to a implementation of ycsb.DB.
type DbWrapper struct {
	DB ycsb.DB
}

func measure(start time.Time, end time.Time, op string, key string, values []interface{}, err error) {
	if err != nil {
		measurement.Measure(fmt.Sprintf("%s_ERROR", op), start, end, key, values)
		return
	}

	measurement.Measure(op, start, end, key, values)
}

func (db DbWrapper) Close() error {
	return db.DB.Close()
}

func (db DbWrapper) InitThread(ctx context.Context, threadID int, threadCount int) context.Context {
	return db.DB.InitThread(ctx, threadID, threadCount)
}

func (db DbWrapper) CleanupThread(ctx context.Context) {
	db.DB.CleanupThread(ctx)
}

func (db DbWrapper) Read(ctx context.Context, table string, key string, fields []string) (_ map[string][]byte, err error) {
	start := time.Now()
	dbRead, err := db.DB.Read(ctx, table, key, fields)
	end := time.Now()
	var tempVals []interface{}
	for _, dbVal := range dbRead {
		tempVals = append(tempVals, dbVal)
	}
	measure(start, end, "READ", key, tempVals, err)

	return dbRead, err
}

func (db DbWrapper) BatchRead(ctx context.Context, table string, keys []string, fields []string) (_ []map[string][]byte, err error) {
	return nil, nil
}

func (db DbWrapper) Scan(ctx context.Context, table string, startKey string, count int, fields []string) (_ []map[string][]byte, err error) {
	return nil, nil
}

func (db DbWrapper) Update(ctx context.Context, table string, key string, values map[string][]byte) (err error) {
	start := time.Now()
	var tempVals []interface{}
	for _, pVal := range values {
		tempVals = append(tempVals, pVal)
	}

	defer func() {
		measure(start, time.Now(), "UPDATE", key, tempVals, err)
	}()

	return db.DB.Update(ctx, table, key, values)
}

func (db DbWrapper) BatchUpdate(ctx context.Context, table string, keys []string, values []map[string][]byte) (err error) {
	return nil
}

func (db DbWrapper) Insert(ctx context.Context, table string, key string, values map[string][]byte) (err error) {
	start := time.Now()
	var tempVals []interface{}
	for _, pVal := range values {
		tempVals = append(tempVals, pVal)
	}

	defer func() {
		measure(start, time.Now(), "INSERT", key, tempVals, err)
	}()

	return db.DB.Insert(ctx, table, key, values)
}

func (db DbWrapper) BatchInsert(ctx context.Context, table string, keys []string, values []map[string][]byte) (err error) {
	return nil
}

func (db DbWrapper) Delete(ctx context.Context, table string, key string) (err error) {
	return nil
}

func (db DbWrapper) BatchDelete(ctx context.Context, table string, keys []string) (err error) {
	return nil
}

func (db DbWrapper) Analyze(ctx context.Context, table string) error {
	return nil
}
