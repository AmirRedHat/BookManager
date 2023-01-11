package utils

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"crypto/sha256"
)


const (
	dbtype = "sqlite3";
	path = "./db.sqlite3";
	write_book_query = "INSERT INTO book VALUES ($1, $2, $3, $4)";
	write_user_query = "INSERT INTO user (username, email, password) VALUES (&1, $2, $3)";
	read_book_query = "SELECT %s FROM book WHERE book_name='%s'";
	read_user_query = "SELECT %s FROM user WHERE id=%d";
	read_all_book_query = "SELECT %s FROM book";
	read_all_user_query = "SELECT %s FROM user";
)

type Store struct {
	path string 
	dbtype string
}

type BookStruct struct {
	Book_name string 
	Author string 
	Views int 
	Timestamp int
}

type UserStruct struct {
	ID int
	Username string 
	Email string 
	Password string
}

func Encrypt(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

func (user *UserStruct) EncryptPassword() {
	current_password := user.Password 
	encrypted_password := Encrypt(current_password)
	user.Password = encrypted_password
}


func (s *Store) StoreWriteBook(book BookStruct) {
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path);
	if dberr != nil{
		log.Fatal(dberr);
	}
	database.Exec(write_book_query, book.Book_name, book.Author, book.Views, book.Timestamp);
	database.Close();
}


func (s *Store) StoreReadBook(book_name string) []BookStruct {
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path);
	if dberr != nil{
		log.Fatal(dberr);
	}
	query := fmt.Sprintf(read_all_book_query, "*")
	if book_name != "all" {
		query = fmt.Sprintf(read_book_query, "*", book_name)
	}
	
	results, dberr := database.Query(query)
	if dberr != nil{
		log.Fatal(dberr);
	}

	book_arr := make([]BookStruct, 0);
	for results.Next(){
		temp_book_struct := BookStruct{}
		scanerr := results.Scan(&temp_book_struct.Book_name, &temp_book_struct.Author, &temp_book_struct.Views, &temp_book_struct.Timestamp)
		if scanerr != nil{
			log.Fatal(scanerr)
		}
		book_arr = append(book_arr, temp_book_struct)
	}

	return book_arr
}

func (s *Store) StoreWriteUser(user UserStruct) {
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path);
	if dberr != nil{
		log.Fatal(dberr);
	}

	// encrypt user password
	user.EncryptPassword()

	database.Exec(write_user_query, user.Username, user.Email, user.Password);
	database.Close();
}

func (s *Store) StoreReadUser(pk int) []UserStruct {
	read_fields_str := "(id, username, email)"
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path);
	if dberr != nil{
		log.Fatal(dberr);
	}
	// creating query
	query := fmt.Sprintf(read_all_user_query, read_fields_str)
	if pk != 0 {
		query = fmt.Sprintf(read_user_query, read_fields_str, pk)
	}
	
	results, dberr := database.Query(query)
	if dberr != nil{
		log.Fatal(dberr);
	}

	user_arr := make([]UserStruct, 0)
	for results.Next() {
		temp_user := UserStruct{}
		results.Scan(&temp_user.ID, &temp_user.Username, &temp_user.Email)
		user_arr = append(user_arr, temp_user)
	}

	return user_arr
}


func WriteBook(book BookStruct) {
	// the middle of saving book
	store := Store{}
	store.StoreWriteBook(book)
}

func ReadBook(book_name string) []BookStruct {
	// the middle of reading book
	store := Store{}
	result := store.StoreReadBook(book_name)
	return result
}