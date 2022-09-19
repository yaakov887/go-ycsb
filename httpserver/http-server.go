package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

var lock = &sync.Mutex{}
var dbLock = &sync.Mutex{}

type fauxDb struct {
	table map[string]int
}

var fauxDbInstance *fauxDb

func getFauxDb() *fauxDb {
	if fauxDbInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if fauxDbInstance == nil {
			fmt.Println("Creating single instance now.")
			fauxDbInstance = &fauxDb{}
			fauxDbInstance.table = make(map[string]int)
		}
	}
	return fauxDbInstance
}

func fauxDbGet(key string) int {
	dbLock.Lock()
	defer dbLock.Unlock()
	fmt.Println("Calling fauxDBGet")

	return getFauxDb().table[key]
}

func fauxDbPut(key string, value int) int {
	dbLock.Lock()
	defer dbLock.Unlock()

	fmt.Println("Calling fauxDBPut")

	getFauxDb().table[key] = value

	return getFauxDb().table[key]
}

func domainParameters(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	fmt.Fprintf(w, "Content Length: %v\n", req.ContentLength)
	fmt.Fprintf(w, "Host: %v\n\n", req.Host)

	fmt.Fprintf(w, "Method: %v\n", req.Method)
	req.ParseForm()
	switch req.Method {
	case "GET":
		for key := range req.Form {
			fmt.Fprintf(w, "Param Key:\t%v\n", key)
			fmt.Fprintf(w, "Returned Value:\t%v\n", fauxDbGet(key))
		}
	case "PUT":
		for key, value := range req.Form {
			tempval, err := strconv.Atoi(value[0])
			if err == nil {
				fmt.Fprintf(w, "Param Key:\t%v\n", key)
				fmt.Fprintf(w, "Put Value:\t%v\n", fauxDbPut(key, tempval))
			}

		}
	}
}

func main() {
	http.HandleFunc("/", domainParameters)

	http.ListenAndServe(":8090", nil)
}
