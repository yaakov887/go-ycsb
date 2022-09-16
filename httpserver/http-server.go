package main

import (
	"fmt"
	"net/http"
)

func headers(w http.ResponseWriter, req *http.Request) {
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
		}
	case "PUT":
		for key, value := range req.Form {
			fmt.Fprintf(w, "Param Key:\t%v\n", key)
			fmt.Fprintf(w, "Param Value:\t%v\n", value)
		}
	}

}

func main() {
	http.HandleFunc("/", headers)

	http.ListenAndServe(":8090", nil)
}
