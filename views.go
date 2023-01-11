package main 

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/thedevsaddam/renderer"
	"localPackage/utils"
	
)

func BookView(res http.ResponseWriter, req *http.Request) {
	rndr := renderer.New()
	method := req.Method

	switch method {
	
	case "GET":
		book_name := req.URL.Query().Get("book_name");
		result := utils.ReadBook(book_name);
		rndr.JSON(res, http.StatusOK, result);
	
	case "POST":
		posted_book := make(map[string]interface{});
		data, err := ioutil.ReadAll(req.Body);
		json.Unmarshal(data, &posted_book)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(posted_book["message"])
	}
}


func User(res http.ResponseWriter, req *http.Request) {
	// rndr := renderer.New()
	method := req.Method 
	switch method {
	
	case "GET":
		user_id := req.URL.Query().Get("pk");
		fmt.Println("user id : ", user_id);
	
	case "POST":
		posted_data, err := ioutil.ReadAll(req.Body);
		if err != nil {
			log.Fatal(err);
		}
		data := make(map[string]interface{})
		json.Unmarshal(posted_data, &data);
		user_name := data["username"];
		user_email := data["email"];
		password := data["password"];
		fmt.Println(user_name, user_email, password);
	}

}