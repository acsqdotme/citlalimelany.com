package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func pageHandler(w http.ResponseWriter, r *http.Request)  {
  if r.URL.Path == "/" {
    r.URL.Path = "/"
  }
  http.ServeFile(w,r,filepath.Join("./html", r.URL.Path))
}

func main()  {
  addr := flag.String("addr", ":4001", "HTTP Server Port Address")
  flag.Parse()

  mux := http.NewServeMux()
  mux.HandleFunc("/", pageHandler)

  log.Printf("Starting server on port %s", *addr)
  http.ListenAndServe(*addr, mux)
}
