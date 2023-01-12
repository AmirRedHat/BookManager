package main

import (
	"encoding/json"
	"fmt"
	"github.com/thedevsaddam/renderer"
	"io/ioutil"
	"localPackage/utils"
	"log"
	"net/http"
	"strconv"
)

func BookView(res http.ResponseWriter, req *http.Request) {
	rndr := renderer.New()
	method := req.Method

	switch method {

	case "GET":
		book_name := req.URL.Query().Get("book_name")
		result := utils.ReadBook(book_name)
		rndr.JSON(res, http.StatusOK, result)

	case "POST":
		posted_book := make(map[string]interface{})
		data, err := ioutil.ReadAll(req.Body)
		json.Unmarshal(data, &posted_book)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(posted_book["message"])
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

		if path == "/register" {
			fmt.Println("registering...")
			posted_data, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			data := make(map[string]interface{})
			json.Unmarshal(posted_data, &data)
			user_name := data["username"].(string)
			user_email := data["email"].(string)
			password := data["password"].(string)
			fmt.Println(user_name, user_email, password)
			//user := utils.UserStruct{Username: user_name, Email: user_email, Password: password}
			//utils.WriteUser(user)

		} else if path == "/login" {

			postedData, _ := ioutil.ReadAll(req.Body)
			data := make(map[string]interface{})
			err := json.Unmarshal(postedData, &data)
			if err != nil {
				log.Fatal(err)
			}
			email := data["email"].(string)
			encryptedPassword := utils.Encrypt(data["password"].(string))
			user := utils.AuthUser(email, encryptedPassword)
			errJson := rndr.JSON(res, http.StatusOK, user)
			if errJson != nil {
				fmt.Fprint(res, errJson)
			}
		}
	}

}
