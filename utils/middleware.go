package utils

import (
	"fmt"
	"github.com/thedevsaddam/renderer"
	"net/http"
)

var rndr = renderer.New()

func AuthTokenMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("middleware is working...")
		email := req.Header.Get("email")
		token := req.Header.Get("token")
		userToken := AuthToken(email, token)
		fmt.Println(userToken)
		if userToken.Email == "" {
			responseMessage := make(map[string]string)
			responseMessage["data"] = "user is not authorized"
			responseMessage["message"] = "failed"
			rndr.JSON(res, http.StatusUnauthorized, responseMessage)
			return
		}
		handler.ServeHTTP(res, req)
	}
}
