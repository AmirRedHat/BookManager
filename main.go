package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	mux := ReturnMux()

	fmt.Println("running server on 0.0.0.0:9090")
	serveerr := http.ListenAndServe("0.0.0.0:9090", mux)
	if serveerr != nil {
		log.Fatal(serveerr)
	}
}
