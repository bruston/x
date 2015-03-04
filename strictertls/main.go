package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

func main() {
	cert := flag.String("cert", "", "Certificate file")
	key := flag.String("key", "", "Private key file")
	listen := flag.String("listen", ":8090", "The host and port to listen on")
	flag.Parse()

	cfg := newRestrictedCfg()
	server := &http.Server{Addr: *listen, TLSConfig: cfg, Handler: http.DefaultServeMux}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})

	log.Printf("starting webserver on %s\n", *listen)

	if err := server.ListenAndServeTLS(*cert, *key); err != nil {
		log.Fatal(err)
	}
}

func newRestrictedCfg() *tls.Config {
	return &tls.Config{
		MinVersion:             tls.VersionTLS12,
		SessionTicketsDisabled: true,
		CipherSuites: []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
