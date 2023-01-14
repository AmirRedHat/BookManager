package utils

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

const (
	dbtype                 = "sqlite3"
	path                   = "./db.sqlite3"
	write_book_query       = "INSERT INTO book VALUES ($1, $2, $3, $4)"
	write_user_query       = "INSERT INTO user (username, email, password) VALUES ($1, $2, $3)"
	write_user_token_query = "INSERT INTO user_token (token, email, time) VALUES (:token, :email, :time)"
	read_book_query        = "SELECT %s FROM book WHERE book_name='%s'"
	read_user_query        = "SELECT %s FROM user WHERE id=%d"
	auth_user_query        = "SELECT %s FROM user WHERE password='%s'"
	read_all_book_query    = "SELECT %s FROM book"
	read_all_user_query    = "SELECT %s FROM user"
)

type Store struct {
	path   string
	dbtype string
}

func (s *Store) returnDB() (*sql.DB, error) {
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path)
	if dberr != nil {
		return nil, dberr
	}
	return database, nil
}

func (s *Store) returnNewDB() *sqlx.DB {
	s.path = path
	s.dbtype = dbtype
	database, err := sqlx.Connect(s.dbtype, s.path)
	if err != nil {
		log.Fatalln(err)
	}

	return database
}

type BookStruct struct {
	BookName  string
	Author    string
	Views     int
	Timestamp int
}

type UserStruct struct {
	ID       int
	Username string
	Email    string
	Password string
}

type UserTokenStruct struct {
	Token      string `db:"token"`
	Email      string `db:"email"`
	ExpireTime int    `db:"time"`
}

func Encrypt(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

// ---------------------------- UserStruct methods
func (user *UserStruct) encryptPassword() {
	currentPassword := user.Password
	encryptedPassword := Encrypt(currentPassword)
	user.Password = encryptedPassword
}

// ---------------------------- UserTokenStruct methods
func (userToken *UserTokenStruct) encryptToken() {
	userToken.Token = Encrypt(fmt.Sprintf("%s%s", userToken.Email, strconv.Itoa(userToken.ExpireTime)))
}

func (s *Store) storeWriteBook(book BookStruct) {
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path)
	if dberr != nil {
		log.Fatal(dberr)
	}
	_, err := database.Exec(write_book_query, book.BookName, book.Author, book.Views, book.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	closeErr := database.Close()
	if closeErr != nil {
		log.Fatal(closeErr)
	}
}

func (s *Store) storeReadBook(book_name string) []BookStruct {
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path)
	if dberr != nil {
		log.Fatal(dberr)
	}
	query := fmt.Sprintf(read_all_book_query, "*")
	if book_name != "all" {
		query = fmt.Sprintf(read_book_query, "*", book_name)
	}

	results, dberr := database.Query(query)
	if dberr != nil {
		log.Fatal(dberr)
	}

	book_arr := make([]BookStruct, 0)
	for results.Next() {
		temp_book_struct := BookStruct{}
		scanerr := results.Scan(&temp_book_struct.BookName, &temp_book_struct.Author, &temp_book_struct.Views, &temp_book_struct.Timestamp)
		if scanerr != nil {
			log.Fatal(scanerr)
		}
		book_arr = append(book_arr, temp_book_struct)
	}

	return book_arr
}

func (s *Store) storeWriteUser(user UserStruct) {
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path)
	if dberr != nil {
		log.Fatal(dberr)
	}
	// encrypt user password
	user.encryptPassword()
	database.Exec(write_user_query, user.Username, user.Email, user.Password)
	database.Close()
}

func (s *Store) storeReadUser(pk int) []UserStruct {
	read_fields_str := "id, username, email"
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path)
	if dberr != nil {
		log.Fatal(dberr)
	}
	// creating query
	query := fmt.Sprintf(read_all_user_query, read_fields_str)
	if pk != 0 {
		query = fmt.Sprintf(read_user_query, read_fields_str, pk)
	}

	results, dberr := database.Query(query)
	if dberr != nil {
		log.Fatal(dberr)
	}

	user_arr := make([]UserStruct, 0)
	for results.Next() {
		temp_user := UserStruct{}
		results.Scan(&temp_user.ID, &temp_user.Username, &temp_user.Email)
		user_arr = append(user_arr, temp_user)
	}

	return user_arr
}

func (s *Store) authPassword(email string, encryptedPassword string) UserStruct {
	read_fields_str := "id, username, email"
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path)
	if dberr != nil {
		log.Fatal(dberr)
	}

	query := fmt.Sprintf(auth_user_query, read_fields_str, encryptedPassword)
	results, dberr := database.Query(query)
	if dberr != nil {
		log.Fatal(dberr)
	}

	targetUser := UserStruct{}
	for results.Next() {
		tempUser := UserStruct{}
		results.Scan(&tempUser.ID, &tempUser.Username, &tempUser.Email)
		if tempUser.Email == email {
			targetUser = tempUser
			break
		}
	}

	return targetUser
}

func (s *Store) StoreWriteUserToken(userToken UserTokenStruct) UserTokenStruct {
	db := s.returnNewDB()
	// encrypt token
	userToken.encryptToken()
	db.NamedExec(write_user_token_query, userToken)
	return userToken
}

func WriteBook(book BookStruct) {
	// the middle of saving book
	store := Store{}
	store.storeWriteBook(book)
}

func ReadBook(book_name string) []BookStruct {
	// the middle of reading book
	store := Store{}
	return store.storeReadBook(book_name)
}

func WriteUser(user UserStruct) {
	store := Store{}
	store.storeWriteUser(user)
}

func ReadUser(pk int) []UserStruct {
	store := Store{}
	return store.storeReadUser(pk)
}

func AuthUser(email string, encryptedPassword string) UserStruct {
	store := Store{}
	return store.authPassword(email, encryptedPassword)
}

func WriteUserToken(userToken UserTokenStruct) UserTokenStruct {
	store := Store{}
	return store.StoreWriteUserToken(userToken)
}
