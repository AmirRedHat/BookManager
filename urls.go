package main

import (
	"localPackage/utils"
	"net/http"
)

func ReturnMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/book", utils.AuthTokenMiddleware(BookView))
	mux.HandleFunc("/user", UserView)
	mux.HandleFunc("/user/register", UserView)
	mux.HandleFunc("/user/login", UserView)
	mux.HandleFunc("/user/logout", UserView)
	//mux.HandleFunc("/user/token", UserTokenView)
	return mux
}
