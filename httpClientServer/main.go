package main

import (
	"flag"
	"fmt"
	"net/http"
)

/*
func head(s string) bool {
	r, e := http.Head(s)
	return e == nil && r.StatusCode == 200
}*/

func handler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "%s   %s from %s \n", req.Method, req.URL.String(), req.Host)
}

func main() {
	//go run main.go -port=8082 for test loadbalancer
	PortNumber := flag.Int(
		"port",
		8081,
		"Port http client server",
	)
	flag.Parse()
	http.HandleFunc("/", handler)
	fmt.Println("Main Server run at 127.0.01 :", *PortNumber)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *PortNumber), nil)
}
