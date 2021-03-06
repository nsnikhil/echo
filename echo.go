package main

import (
	"log"
	"net/http"
	"os"
)

func startServer(address string) {
	writeData := func(code int, msg string, resp http.ResponseWriter) {
		resp.WriteHeader(code)
		_, _ = resp.Write([]byte(msg))
	}

	http.HandleFunc("/echo", func(resp http.ResponseWriter, req *http.Request) {
		val, ok := req.URL.Query()["value"]
		if !ok || len(val) == 0 || len(val[0]) == 0 {
			log.Printf("invalid url, value is missing: %s", req.URL.String())
			writeData(http.StatusBadRequest, "invalid url, value is missing", resp)
			return
		}

		log.Printf("echo success: %s", val[0])
		writeData(http.StatusOK, val[0], resp)
	})

	log.Printf("listening on: %s", address)
	_ = http.ListenAndServe(address, nil)
}

func main() {
	startServer(os.Getenv("ADDRESS"))
}
