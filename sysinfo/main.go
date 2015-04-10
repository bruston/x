package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:5432", "The host and port to listen on.")
	flag.Parse()
	http.HandleFunc("/env", env)
	http.HandleFunc("/uptime", uptime)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

func env(w http.ResponseWriter, r *http.Request) {
	envDump := os.Environ()
	envMap := make(map[string]string)
	for _, v := range envDump {
		pair := strings.SplitN(v, "=", 2)
		envMap[pair[0]] = pair[1]
	}
	writeJSON(w, envMap, http.StatusOK)
}

func uptime(w http.ResponseWriter, r *http.Request) {
	si := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(si); err != nil {
		writeJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uptime := time.Duration(si.Uptime) * time.Second
	writeJSON(w, uptime.String(), http.StatusOK)
}
