package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var StaticFiles embed.FS

func NewHandler() http.Handler {
	mux := http.NewServeMux()

	staticFS, _ := fs.Sub(StaticFiles, "static")

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
			return
		}
		data, err := StaticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	return mux
}
