package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4001", "HTTP Server Port Address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", pageHandler)
	mux.HandleFunc("/favicon.ico", faviconHandler)

	log.Printf("Starting server on port %s", *addr)
	http.ListenAndServe(*addr, mux)
}
