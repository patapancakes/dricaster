package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/netip"

	"github.com/patapancakes/sslspoof"
)

func main() {
	fmt.Println("Dricaster by Pancakes (pancakes@mooglepowered.com)")
	fmt.Println()

	addr := flag.String("addr", "0.0.0.0:443", "address to listen on")
	flag.Parse()

	l, err := sslspoof.NewListener(*addr, "auth01.dricas.com", true)
	if err != nil {
		panic(err)
	}

	defer l.Close()

	http.HandleFunc("POST /cgi-bin/auth.cgi", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		addrport, _ := netip.ParseAddrPort(r.RemoteAddr)
		log.Printf("[%s] %v", addrport.Addr(), r.PostForm)
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Listening on", *addr)

	err = http.Serve(l, nil)
	if err != nil {
		panic(err)
	}
}
