package main 

import (
	"localPackage/utils"
	"net/http"
	"log"
	"fmt"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/book", BookView)
	
	fmt.Println("running server on 0.0.0.0:9090")
	serveerr := http.ListenAndServe("0.0.0.0:9090", mux)
	if serveerr != nil {
		log.Fatal(serveerr)
	}

	fmt.Println(utils.Entcrypt("python"))
	// result := utils.ReadBook("bookName1")
}