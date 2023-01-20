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
		result, err := utils.ReadBook(book_name)
		if err != nil {
			responseMessage := make(map[string]interface{})
			responseMessage["error"] = err.Error()
			rndr.JSON(res, http.StatusInternalServerError, responseMessage)
		}
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
			responseData := make(map[string]interface{})
			user := utils.UserStruct{Username: userName, Email: userEmail, Password: password}
			utils.WriteUser(user)

			exTime := int(time.Now().Add(30 * time.Minute).Unix())
			token := ""
			userToken := utils.UserTokenStruct{Token: token, Email: userEmail, ExpireTime: exTime}
			userToken = utils.WriteUserToken(userToken)

			responseData["user"] = user
			responseData["token"] = userToken
			rndr.JSON(res, http.StatusCreated, responseData)

		} else if path == "/user/login" {

			//loginTime := time.Now()
			data := returnData(req)
			email := data["email"].(string)
			password := data["password"].(string)

			user, err := utils.AuthPassword(email, password)
			if user.Email == "" || err != nil {
				rndr.JSON(res, http.StatusUnauthorized, user)
				return
			}

			exTime := int(time.Now().Add(30 * time.Minute).Unix())
			token := ""
			userToken := utils.UserTokenStruct{Token: token, Email: email, ExpireTime: exTime}
			//userTokenCreationTime := time.Now()
			// removing all active tokens with this email
			utils.DestroyToken(email)
			// write new user token in db
			userToken = utils.WriteUserToken(userToken)
			//fmt.Println("duration in creating token : ", time.Since(userTokenCreationTime).Seconds())
			rndr.JSON(res, http.StatusOK, userToken)
			//fmt.Println("duration in logging in : ", time.Since(loginTime).Seconds())

		} else if path == "/user/logout" {
			data := returnData(req)
			email := data["email"].(string)
			utils.DestroyToken(email)
			rndr.JSON(res, http.StatusOK, "user logged out")
		}
	}

}

//func UserTokenView(res http.ResponseWriter, req *http.Request) {
//	rndr := renderer.New()
//	method := req.Method
//	path := req.URL.Path
//	data := returnData(req)
//
//	if method == "POST" {
//		if path == "/user/token/access" || path == "/user/token/" {
//			email := data["email"].(string)
//			exTime := int(time.Now().Add(30 * time.Minute).Unix())
//			token := ""
//			userToken := utils.UserTokenStruct{Token: token, Email: email, ExpireTime: exTime}
//			userToken = utils.WriteUserToken(userToken)
//
//		} else if path == "/user/token/auth" {
//			email := data["email"].(string)
//			token := data["token"].(string)
//			userToken := utils.AuthToken(email, token)
//			status := http.StatusOK
//			if userToken.Email == "" {
//				status = http.StatusUnauthorized
//			}
//			rndr.JSON(res, status, userToken)
//		}
//	}
//}
