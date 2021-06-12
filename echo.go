package main

import (
	"log"
	"net/http"
	"os"
)

func echoHandler(resp http.ResponseWriter, req *http.Request) {
	val, ok := req.URL.Query()["value"]

	if !ok || len(val) == 0 || len(val[0]) == 0 {
		log.Printf("invalid url, value is missing: %s", req.URL.String())
		resp.WriteHeader(http.StatusBadRequest)
		_, _ = resp.Write([]byte("invalid url, value is missing"))
		return
	}

	log.Printf("echo success: %s", val[0])
	resp.WriteHeader(http.StatusOK)
	_, _ = resp.Write([]byte(val[0]))
}

func startServer(address string) {
	http.HandleFunc("/echo", echoHandler)

	log.Printf("listening on: %s", address)
	_ = http.ListenAndServe(address, nil)
}

func main() {
	address := os.Getenv("ADDRESS")
	if len(address) == 0 {
		log.Fatal("please provide a valid address")
	}

	startServer(address)
}
