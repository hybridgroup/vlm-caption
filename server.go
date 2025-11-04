package main

import (
	"log"
	"net/http"
	"time"

	_ "embed"

	"github.com/hybridgroup/mjpeg"
)

//go:embed html/index.html
var index string

func startWebServer(host string, stream *mjpeg.Stream) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(index))
	})
	mux.Handle("/video", stream)
	mux.HandleFunc("/caption", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(caption))
	})

	server := &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
