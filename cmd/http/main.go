package main

import (
	"flag"
	"log"
	"net/http"

	"citlalimelany.com/album"
)

func main() {
	if err := album.MakeDB(); err != nil {
		log.Fatal(err)
	}
	addr := flag.String("addr", ":4001", "HTTP Server Port Address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", pageHandler)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/favicon.ico", faviconHandler)

	log.Printf("Starting server on port %s", *addr)
	http.ListenAndServe(*addr, mux)
}
