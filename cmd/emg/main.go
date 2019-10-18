package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/AgentZombie/emg"
)

var (
	fListen = flag.String("listen", ":8080", "listen on this address/port")
)

func fatalIfError(err error, msg string) {
	if err != nil {
		log.Fatal("error ", msg, ": ", err)
	}
}

func main() {
	flag.Parse()
	_, err := emg.New()
	fatalIfError(err, "creating server")
	log.Print("starting listener on ", *fListen)
	fatalIfError(http.ListenAndServe(*fListen, nil), "listening")
}
