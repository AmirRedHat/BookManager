package main

import (
	"fmt"
	"localPackage/utils"
	"log"
	"net/http"
)

func controller() {
	mux := http.NewServeMux()
	mux.HandleFunc("/book", BookView)
	mux.HandleFunc("/user", UserView)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/book", BookView)
	mux.HandleFunc("/user", UserView)

	fmt.Println("running server on 0.0.0.0:9090")
	serveerr := http.ListenAndServe("0.0.0.0:9090", mux)
	if serveerr != nil {
		log.Fatal(serveerr)
	}

	fmt.Println(utils.Entcrypt("python"))
	result := utils.ReadUser(1)
	fmt.Println(result)
}
