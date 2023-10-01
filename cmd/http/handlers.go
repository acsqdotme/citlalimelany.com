package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"citlalimelany.com/album"
)

const (
	htmlDir     = "./html"
	staticDir   = "./static"
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

	data := make(map[string]interface{})
	data["Path"] = r.URL.Path
	albums, err := album.AggregateAlbums()
	if err != nil {
		log.Fatal(err.Error())
	}
	data["Albums"] = albums
	data["Album"] = albums[0]

	tmpl.ExecuteTemplate(w, "base", data)
}

func albumHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	path := strings.Split(r.URL.Path, "/")
	pageAlbum := path[2]

	if len(path) > 4 || (len(path) == 4 && path[3] != "") {
		http.Error(w, "page not found cause len", http.StatusNotFound)
		return
	} else if len(path) == 4 && path[3] == "" {
		http.Redirect(w, r, scheme+"://"+r.Host+"/album/"+pageAlbum, 302)
		return
	} else if len(path) == 3 && pageAlbum == "" {
		http.Redirect(w, r, scheme+"://"+r.Host+"/", 302)
		return
	}

	exists, err := album.DoesAlbumExist(pageAlbum)

	if !exists {
		http.Error(w, "page not found cause bool", 404)
		return
	}

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "album"+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	data := make(map[string]interface{})
	data["Path"] = r.URL.Path
	a, err := album.FetchAlbum(pageAlbum)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}
	data["Album"] = a

	tmpl.ExecuteTemplate(w, "base", data)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir, "favicon.ico"))
}
