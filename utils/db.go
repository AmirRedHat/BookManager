package utils

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	dbtype                  = "sqlite3"
	path                    = "./db.sqlite3"
	write_book_query        = "INSERT INTO book VALUES ($1, $2, $3, $4)"
	write_user_query        = "INSERT INTO user (username, email, password) VALUES ($1, $2, $3)"
	write_user_token_query  = "INSERT INTO user_token VALUES ($1, $2, $3)"
	delete_user_token_query = "DELETE FROM user_token WHERE email='%s'"
	read_user_token_query   = "SELECT %s FROM user_token WHERE token='%s' AND email='%s'"
	read_book_query         = "SELECT %s FROM book WHERE book_name='%s'"
	read_user_query         = "SELECT %s FROM user WHERE id=%d"
	auth_user_query         = "SELECT %s FROM user WHERE email='%s' AND password='%s'"
	read_all_book_query     = "SELECT %s FROM book"
	read_all_user_query     = "SELECT %s FROM user"
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
func (book *BookStruct) validateBook() bool {
	if strings.Contains(book.BookName, ";") || strings.Contains(book.BookName, "'") {
		return false
	}
	return true
}

// ---------------------------- UserStruct methods
func (user *UserStruct) encryptPassword() {
	currentPassword := user.Password
	encryptedPassword := Encrypt(currentPassword)
	user.Password = encryptedPassword
}

func (user *UserStruct) validateUser() bool {
	emailCondition := !strings.Contains(user.Email, "@gmail.com")
	validateLetters := strings.Contains(user.Email, "'") || strings.Contains(user.Email, ";")
	if validateLetters && emailCondition {
		return false
	}
	return true
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

func (s *Store) storeReadBook(book_name string) ([]BookStruct, error) {
	s.path = path
	s.dbtype = dbtype
	book_arr := make([]BookStruct, 0)
	database, dberr := s.returnDB()
	if dberr != nil {
		log.Fatal(dberr)
	}

	query := fmt.Sprintf(read_all_book_query, "*")
	if book_name != "all" {
		query = fmt.Sprintf(read_book_query, "*", book_name)
	}

	results, dberr := database.Query(query)
	if dberr != nil {
		return book_arr, dberr
	}

	for results.Next() {
		temp_book_struct := BookStruct{}
		scanerr := results.Scan(&temp_book_struct.BookName, &temp_book_struct.Author, &temp_book_struct.Views, &temp_book_struct.Timestamp)
		if scanerr != nil {
			log.Fatal(scanerr)
		}
		book_arr = append(book_arr, temp_book_struct)
	}

	database.Close()
	results.Close()
	return book_arr, nil
}

func (s *Store) storeWriteUser(user UserStruct) {
	s.path = path
	s.dbtype = dbtype
	database, dberr := sql.Open(s.dbtype, s.path)
	if dberr != nil {
		fmt.Println("open storeWriteUser")
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
		fmt.Println("open storeReadUser")
		log.Fatal(dberr)
	}
	// creating query
	query := fmt.Sprintf(read_all_user_query, read_fields_str)
	if pk != 0 {
		query = fmt.Sprintf(read_user_query, read_fields_str, pk)
	}

	results, dberr := database.Query(query)
	if dberr != nil {
		fmt.Println("query storeReadUser")
		log.Fatal(dberr)
	}

	user_arr := make([]UserStruct, 0)
	for results.Next() {
		temp_user := UserStruct{}
		results.Scan(&temp_user.ID, &temp_user.Username, &temp_user.Email)
		user_arr = append(user_arr, temp_user)
	}

	database.Close()
	results.Close()
	return user_arr
}

func (s *Store) authPassword(email string, password string) (UserStruct, error) {
	read_fields_str := "id, username, email"
	s.path = path
	s.dbtype = dbtype
	targetUser := UserStruct{}
	database, dberr := sql.Open(s.dbtype, s.path)
	if dberr != nil {
		fmt.Println("authPassword")
		log.Fatal(dberr)
	}

	//queryTime := time.Now()
	encryptedPassword := Encrypt(password)
	query := fmt.Sprintf(auth_user_query, read_fields_str, email, encryptedPassword)
	results, dberr := database.Query(query)
	if dberr != nil {
		fmt.Println("query error: ", dberr)
		return targetUser, dberr
	}
	//fmt.Println("duration in encrypting and querying : ", time.Since(queryTime).Seconds())

	//scanTime := time.Now()
	results.Next()
	results.Scan(&targetUser.ID, &targetUser.Username, &targetUser.Email)
	//fmt.Println("duration in scanning : ", time.Since(scanTime).Seconds())

	database.Close()
	results.Close()

	return targetUser, nil
}

func (s *Store) storeWriteUserToken(userToken UserTokenStruct) UserTokenStruct {
	db, err := s.returnDB()
	if err != nil {
		log.Fatal(err)
	}
	// encrypt token
	writeTokenTime := time.Now()
	userToken.encryptToken()
	fmt.Println("duration in encrypting token: ", time.Since(writeTokenTime).Seconds())
	results, err := db.Exec(write_user_token_query, userToken.Token, userToken.Email, userToken.ExpireTime)
	if err != nil {
		fmt.Println(results)
		log.Fatal(err)
	}
	db.Close()
	fmt.Println("duration in write token : ", time.Since(writeTokenTime).Seconds())
	return userToken
}

func (s *Store) authUserToken(email string, token string) UserTokenStruct {
	db := s.returnNewDB()
	query := fmt.Sprintf(read_user_token_query, "*", token, email)
	results, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	} else {
		db.Close()
	}

	count := 0
	var UserToken UserTokenStruct
	for results.Next() {
		tempUserToken := UserTokenStruct{}
		if count > 0 {
			log.Fatal("internal error : too many token")
		}
		count++
		results.Scan(&tempUserToken.Token, &tempUserToken.Email, &tempUserToken.ExpireTime)
		UserToken = tempUserToken
	}

	nowTimestamp := int(time.Now().Unix())
	if nowTimestamp > UserToken.ExpireTime {
		return UserTokenStruct{}
	}
	return UserToken
}

func (s *Store) destroyToken(email string) {
	db := s.returnNewDB()
	query := fmt.Sprintf(delete_user_token_query, email)
	fmt.Println("delete : ", query)
	result, err := db.Exec(query)
	if err != nil {
		fmt.Println(result)
		fmt.Println("destroy token error : ", err)
		log.Fatal(err)
	}
	db.Close()
}

var store Store

func WriteBook(book BookStruct) {
	// the middle of saving book

	store.storeWriteBook(book)
}

func ReadBook(bookName string) ([]BookStruct, error) {
	// the middle of reading book
	if strings.Contains(bookName, "'") || strings.Contains(bookName, ";") {
		return []BookStruct{}, errors.New("invalid book name (cannot contain ' or ;)")
	}
	return store.storeReadBook(bookName)

}

func WriteUser(user UserStruct) {
	store.storeWriteUser(user)
}

func ReadUser(pk int) []UserStruct {
	return store.storeReadUser(pk)
}

func AuthPassword(email string, encryptedPassword string) (UserStruct, error) {
	return store.authPassword(email, encryptedPassword)
}

func WriteUserToken(userToken UserTokenStruct) UserTokenStruct {
	return store.storeWriteUserToken(userToken)
}

func AuthToken(email string, token string) UserTokenStruct {
	return store.authUserToken(email, token)
}

func DestroyToken(email string) {
	store.destroyToken(email)
}
