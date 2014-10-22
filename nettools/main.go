package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	var listen = flag.String("listen", ":8090", `The host and port to listen on. Examples: ":8090" listens on all interfaces on port 8090. "127.0.0.1:8090" listens on the loopback interface on port 8090. Defaults to: ":8090".`)
	flag.Parse()
	m := mux.NewRouter()
	m.HandleFunc("/ip", getIP)
	m.HandleFunc("/port/{port}", tryPort)
	m.HandleFunc("/hostnames", getHostnames)
	log.Fatal(http.ListenAndServe(*listen, m))
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

type jError struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

func JSONError(w http.ResponseWriter, message string, status int) error {
	return writeJSON(w, jError{message, status}, status)
}

type IP struct {
	IP string `json:"ip"`
}

func addrFromRequest(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func getIP(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, IP{IP: addrFromRequest(r)}, http.StatusOK)
}

type portStatus struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Open  bool   `json:"open"`
	Error string `json:"error"`
}

func tryPort(w http.ResponseWriter, r *http.Request) {
	port, err := strconv.Atoi(mux.Vars(r)["port"])
	if err != nil {
		JSONError(w, "bad request, supplied port must be a number", http.StatusBadRequest)
		return
	}
	ip := addrFromRequest(r)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		writeJSON(w, portStatus{ip, port, false, err.Error()}, http.StatusOK)
		return
	}
	conn.Close()
	writeJSON(w, portStatus{ip, port, true, ""}, http.StatusOK)
}

type hosts struct {
	Hostnames []string `json:"hostnames"`
}

func getHostnames(w http.ResponseWriter, r *http.Request) {
	ip := addrFromRequest(r)
	var results hosts
	var err error
	results.Hostnames, err = net.LookupAddr(ip)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, results, http.StatusOK)
}
