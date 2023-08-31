package main

import (
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	htmlDir     = "./html"
	tmplFileExt = ".tmpl.html"
)

func doesFileExist(pathToFile string) bool {
	info, err := os.Stat(filepath.Clean(pathToFile))
	if err != nil || info.IsDir() {
		return false
	}
	return true
}

func bindTMPL(files ...string) (*template.Template, error) {
	for _, checkFile := range files {
		if !doesFileExist(checkFile) {
			return nil, errors.New("Template file missing " + checkFile)
		}
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	path := strings.Split(r.URL.Path, "/")
	page := path[1]

	if r.URL.Path == "/" {
		page = "index"
	}

	if !doesFileExist(filepath.Join(htmlDir, "pages", page+tmplFileExt)) {
		http.Error(w, "page not found", 404)
		return
	}

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "pages", page+tmplFileExt),
	)
	if err != nil {
		http.Error(w, "template broke", 500)
		log.Println(err.Error())
		return
	}

	tmpl.ExecuteTemplate(w, "base", nil)
}

func main() {
	addr := flag.String("addr", ":4001", "HTTP Server Port Address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", pageHandler)

	log.Printf("Starting server on port %s", *addr)
	http.ListenAndServe(*addr, mux)
}
