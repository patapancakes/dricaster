package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	fmt.Println("Dricaster by Pancakes (pancakes@mooglepowered.com)")
	fmt.Println()

	addr := flag.String("addr", "0.0.0.0:443", "address to listen on")
	flag.Parse()

	err := setupSSL()
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	listener := sslListener{Listener: l}
	defer listener.Close()

	http.HandleFunc("POST /cgi-bin/auth.cgi", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got auth request from", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Listening on", *addr)

	err = http.Serve(&listener, nil)
	if err != nil {
		panic(err)
	}
}
