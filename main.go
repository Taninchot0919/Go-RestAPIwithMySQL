package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Model
type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Func Main
func main() {
	fmt.Println("Start project")
	db, err = sql.Open("mysql", "root:mysql@tcp(168.63.232.208:3306)/PostDB") // ต้องใช้ = เท่านั้น ห้ามใช้ := เพราะเราประกาศไว้ข้างบนแล้ว หากใช้ := จะเกิด memory error

	if err != nil {
		panic(err)
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")

	http.ListenAndServe(":9000", router)
}

// Below is method for REST API
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var posts []Post
	result, err := db.Query("SELECT id, title, body from Post")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var post Post
		err := result.Scan(&post.ID, &post.Title, &post.Body)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	fmt.Println("GetAll Method Run Successfully!")

	json.NewEncoder(w).Encode(posts)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)

	result, err := db.Query("SELECT id,title,body FROM Post where id = ?", param["id"])
	if err != nil {
		panic(err)
	}

	defer result.Close()

	var post Post
	for result.Next() {
		err := result.Scan(&post.ID, &post.Title, &post.Body)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("GetOne Method Run Successfully!")

	json.NewEncoder(w).Encode(post)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	statement, err := db.Prepare("INSERT INTO Post(title,body) VALUE(?,?)")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	keyValue := make(map[string]string)
	json.Unmarshal(body, &keyValue)
	title := keyValue["title"]
	bodyData := keyValue["body"]

	_, err = statement.Exec(title, bodyData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Post Method Run Successfully!")
	fmt.Fprintf(w, "New post was create")
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	statement, err := db.Prepare("UPDATE Post SET title = ? , body = ? where id = ?")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	keyValue := make(map[string]string)
	json.Unmarshal(body, &keyValue)
	newTitle := keyValue["title"]
	newBody := keyValue["body"]

	_, err = statement.Exec(newTitle, newBody, param["id"])
	if err != nil {
		panic(err)
	}

	fmt.Println("Update Method Run Successfully!")
	fmt.Fprintf(w, "Post with ID = %s was updated", param["id"])
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	statement, err := db.Prepare("DELETE FROM Post WHERE id = ?")
	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(param["id"])
	if err != nil {
		panic(err)
	}

	fmt.Println("Delete Method Run Successfully!")
	fmt.Fprintf(w, "Delete ID : %s successfully", param["id"])
}
