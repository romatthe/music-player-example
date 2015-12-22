package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

// FileInfo is a nice little struct
type FileInfo struct {
	Name  string
	IsDir bool
	Mode  os.FileMode
}

const (
	filePrefix = "/music/"
	root       = "./music"
)

func main() {
	http.HandleFunc("/", playerMainFrame)
	http.HandleFunc(filePrefix, serveFile)
	http.ListenAndServe(":8080", nil)
}

func playerMainFrame(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "player.html")
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(root, r.URL.Path[len(filePrefix):])
	stat, err := os.Stat(path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	if stat.IsDir() {
		serveDir(w, r, path)
		return
	}

	http.ServeFile(w, r, path)
}

func serveDir(w http.ResponseWriter, r *http.Request, path string) {
	defer func() {
		if err, ok := recover().(error); ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	files, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fileInfos := make([]FileInfo, len(files), len(files))

	for i := range files {
		fileInfos[i].Name = files[i].Name()
		fileInfos[i].Mode = files[i].Mode()
		fileInfos[i].IsDir = files[i].IsDir()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&fileInfos); err != nil {
		panic(err)
	}
}
