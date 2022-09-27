package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var lock = &sync.Mutex{}
var dbLock = &sync.Mutex{}

type fauxDb struct {
	table map[string]string
}

var fauxDbInstance *fauxDb

func getFauxDb() *fauxDb {
	if fauxDbInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if fauxDbInstance == nil {
			fmt.Println("Creating single instance now.")
			fauxDbInstance = &fauxDb{}
			fauxDbInstance.table = make(map[string]string)
		}
	}
	return fauxDbInstance
}

func fauxDbGet(key string) string {
	dbLock.Lock()
	defer dbLock.Unlock()

	value := getFauxDb().table[key]
	fmt.Printf("Calling fauxDB GET - %v:%v\n", key, value)

	return value
}

func fauxDbPut(key string, value string) string {
	dbLock.Lock()
	defer dbLock.Unlock()

	fmt.Printf("Calling fauxDB PUT - %v:%v\n", key, value)

	getFauxDb().table[key] = value

	return getFauxDb().table[key]
}

func domainParameters(w http.ResponseWriter, req *http.Request) {
	//debug lines to print the incoming request parameters
	//for name, headers := range req.Header {
	//	for _, h := range headers {
	//		fmt.Printf("%v: %v\n", name, h)
	//	}
	//}
	//fmt.Printf("Content Length: %v\n", req.ContentLength)
	//fmt.Printf("Host: %v\n\n", req.Host)
	//fmt.Printf("Method: %v\n", req.Method)

	fmt.Printf("Method : %v\n", req.Method)
	switch req.Method {
	case "GET":
		fmt.Printf("Path : %+v\n", req.URL.Path)
		paths := strings.Split(req.URL.Path, "/")
		var key string
		if paths[0] == "" {
			key = paths[1]
		} else {
			key = paths[0]
		}
		fmt.Printf("key: %v\n", key)
		fmt.Fprintf(w, "%v", fauxDbGet(key))
	case "PUT", "POST":
		key := strings.Split(req.URL.Path, "/")[1]
		fmt.Printf("key: %v\n", key)
		value := req.URL.Query()
		fmt.Printf("value: %+v\n", value)
		fmt.Fprintf(w, "%v:%v,", key, fauxDbPut(key, value["value"][0]))
	}
}

func main() {
	http.HandleFunc("/", domainParameters)

	http.ListenAndServe(":8090", nil)
}
