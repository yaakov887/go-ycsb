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

package httpdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/pingcap/go-ycsb/pkg/ycsb"
	"io"
	"net/http"
)

// http properties
const (
	httpDomain = "domain"
	httpPort   = "port"
)

type httpDB struct {
	p      *properties.Properties
	domain string
	port   string
	conn   *http.Client
}

type httpDBCreator struct{}

func (c httpDBCreator) Create(p *properties.Properties) (ycsb.DB, error) {
	db := new(httpDB)

	db.p = p
	db.domain = p.GetString(httpDomain, "localhost")
	db.port = p.GetString(httpPort, "8090")

	db.conn = &http.Client{}
	fmt.Printf("%+v\n", db)

	return db, nil
}

func (h httpDB) GenerateURL(key string) string {
	if h.port == "" {
		return fmt.Sprintf("http://%v/%v", h.domain, key)
	} else {
		return fmt.Sprintf("http://%v:%v/%v", h.domain, h.port, key)
	}
}

// Close closes the database layer.
func (h httpDB) Close() error {
	return nil
}

// InitThread initializes the state associated to the goroutine worker.
// The Returned context will be passed to the following usage.
func (h httpDB) InitThread(ctx context.Context, threadID int, threadCount int) context.Context {
	return ctx
}

// CleanupThread cleans up the state when the worker finished.
func (h httpDB) CleanupThread(ctx context.Context) {}

// Read reads a record from the database and returns a map of each field/value pair.
// table: The name of the table.
// key: The record key of the record to read.
// fields: The list of fields to read, nil|empty for reading all.
func (h httpDB) Read(ctx context.Context, table string, key string, fields []string) (map[string][]byte, error) {
	tempURL := h.GenerateURL(key)
	//fmt.Printf("GET tempURL: %v\n", tempURL)

	resp, err := h.conn.Get(tempURL)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		//println(string(body))
		readVal := make(map[string][]byte)
		readVal[fields[0]] = body
		return readVal, nil
	}
	return nil, err
}

// Scan scans records from the database.
// table: The name of the table.
// startKey: The first record key to read.
// count: The number of records to read.
// fields: The list of fields to read, nil|empty for reading all.
func (h httpDB) Scan(ctx context.Context, table string, startKey string, count int, fields []string) ([]map[string][]byte, error) {
	return nil, errors.New("scan not implemented")
}

// Update updates a record in the database. Any field/value pairs will be written into the
// database or overwritten the existing values with the same field name.
// table: The name of the table.
// key: The record key of the record to update.
// values: A map of field/value pairs to update in the record.
func (h httpDB) Update(ctx context.Context, table string, key string, values map[string][]byte) error {
	var tempVal string
	for _, v := range values {
		tempVal = string(v)
		break
	}
	tempURL := h.GenerateURL(key)
	tempURL = fmt.Sprintf("%v?value=%v", tempURL, tempVal)
	//fmt.Printf("PUT tempURL: %v\n", tempURL)

	req, err := http.NewRequest("PUT", tempURL, nil)
	_, err = h.conn.Do(req)

	if err != nil {
		return err
	}
	//println(resp.Body)
	return nil
}

// Insert inserts a record in the database. Any field/value pairs will be written into the
// database.
// table: The name of the table.
// key: The record key of the record to insert.
// values: A map of field/value pairs to insert in the record.
func (h httpDB) Insert(ctx context.Context, table string, key string, values map[string][]byte) error {
	var tempVal string
	for _, v := range values {
		tempVal = string(v)
		break
	}
	tempURL := h.GenerateURL(key)
	tempURL = fmt.Sprintf("%v?value=%v", tempURL, tempVal)
	//fmt.Printf("PUT tempURL: %v\n", tempURL)

	req, err := http.NewRequest("PUT", tempURL, nil)
	_, err = h.conn.Do(req)

	if err != nil {
		return err
	}
	//println(resp.Body)
	return nil
}

// Delete deletes a record from the database.
// table: The name of the table.
// key: The record key of the record to delete.
func (h httpDB) Delete(ctx context.Context, table string, key string) error {
	return errors.New("delete not implemented")
}

func init() {
	fmt.Println("Registering httpdb...")
	ycsb.RegisterDBCreator("httpdb", httpDBCreator{})
}
