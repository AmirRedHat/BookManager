package main

import (
	"encoding/json"
	"github.com/thedevsaddam/renderer"
	"io/ioutil"
	"localPackage/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

func returnData(req *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	posted_data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(posted_data, &data)
	return data
}

func BookView(res http.ResponseWriter, req *http.Request) {
	rndr := renderer.New()
	method := req.Method

	switch method {

	case "GET":
		book_name := req.URL.Query().Get("book_name")
		result := utils.ReadBook(book_name)
		rndr.JSON(res, http.StatusOK, result)

	case "POST":
		data := make(map[string]interface{})
		posted_book, err := ioutil.ReadAll(req.Body)
		json.Unmarshal(posted_book, &data)
		if err != nil {
			log.Fatal(err)
		}
		bookName := data["book_name"].(string)
		author := data["author"].(string)
		book := utils.BookStruct{BookName: bookName, Author: author, Views: 0, Timestamp: int(time.Now().Unix())}
		utils.WriteBook(book)
		rndr.JSON(res, http.StatusCreated, book)
	}
}

func UserView(res http.ResponseWriter, req *http.Request) {
	rndr := renderer.New()
	method := req.Method
	path := req.URL.Path

	switch method {
	case "GET":
		user_id := req.URL.Query().Get("pk")
		id, err := strconv.Atoi(user_id)
		if err != nil {
			log.Fatal(err)
		}
		result := utils.ReadUser(id)
		rndr.JSON(res, http.StatusOK, result)

	case "POST":

		if path == "/user/register" {
			data := returnData(req)
			userName := data["username"].(string)
			userEmail := data["email"].(string)
			password := data["password"].(string)
			user := utils.UserStruct{Username: userName, Email: userEmail, Password: password}
			utils.WriteUser(user)
			rndr.JSON(res, http.StatusCreated, user)

		} else if path == "/user/login" {

			data := returnData(req)
			email := data["email"].(string)
			encryptedPassword := utils.Encrypt(data["password"].(string))
			user := utils.AuthUser(email, encryptedPassword)
			rndr.JSON(res, http.StatusOK, user)

		} else if path == "/user/logout" {
			data := returnData(req)
			email := data["email"].(string)
			utils.DestroyToken(email)
			rndr.JSON(res, http.StatusOK, "user logged out")
		}
	}

}

func UserTokenView(res http.ResponseWriter, req *http.Request) {
	rndr := renderer.New()
	method := req.Method
	path := req.URL.Path
	data := returnData(req)

	if method == "POST" {
		if path == "/user/token/access" || path == "/user/token/" {
			email := data["email"].(string)
			time := int(time.Now().Unix())
			token := ""
			userToken := utils.UserTokenStruct{Token: token, Email: email, ExpireTime: time}
			userToken = utils.WriteUserToken(userToken)

		} else if path == "/user/token/auth" {
			email := data["email"].(string)
			token := data["token"].(string)
			userToken := utils.AuthToken(email, token)
			status := http.StatusOK
			if userToken.Email == "" {
				status = http.StatusUnauthorized
			}
			rndr.JSON(res, status, userToken)
		}
	}
}
