package main

import (
	"flag"
	"log"
	"net/http"

	"citlalimelany.com/album"
)

var (
	scheme = "http"
)

func main() {
	if err := album.MakeDB(); err != nil {
		log.Fatal(err)
	}

	addr := flag.String("addr", ":4001", "HTTP Server Port Address")
	https := flag.Bool("https", false, "Enable HTTPS for Deployment")
	flag.Parse()

	if *https {
		scheme = "https"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", pageHandler)
	mux.HandleFunc("/album/", albumHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/favicon.ico", faviconHandler)

	log.Printf("Starting server on port %s", *addr)
	http.ListenAndServe(*addr, mux)
}
