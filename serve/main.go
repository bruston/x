package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/bruston/handlers/debug"
)

func main() {
	listen := flag.String("listen", "localhost:8000", "The host:port to listen on.")
	dir := flag.String("dir", "", "The directory to serve. Current directory used if left empty.")
	showDebug := flag.Bool("debug", true, "Log request information to the console.")
	flag.Parse()

	var handler http.Handler = http.DefaultServeMux
	if *showDebug {
		handler = debug.New(handler)
	}

	if *dir == "" {
		d, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dir = &d
	}

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	log.Printf("Webserver starting on: http://%s", *listen)
	log.Fatal(http.ListenAndServe(*listen, handler))
}
