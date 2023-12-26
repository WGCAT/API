package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct { //제이슨이 읽을 수 있는 유저 스트럭트를 만듦
	ID        int       `json:"ID"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}
func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get UserInfo by /users/{id}")
}
func getUserInfoHandler(w http.ResponseWriter, r *http.Request) { //고정 89가 아니라 아이디를 나타내줘야하므로 mux.Vars사용한다
	user := new(User)
	user.ID = 2
	user.FirstName = "sujin"
	user.LastName = "lee"
	user.Email = "seed9878@gmail.com"

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}
func CreateUserHandler(w http.ResponseWriter, r *http.Request) { //실제 유저를 생성하는 코드를 만들어야하는데 클라이언트가 유저정보를 제이슨으로 보냈음
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	// Created User
	user.ID = 2
	user.CreatedAt = time.Now()
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

// NewHandler make a new myapp handler
func NewHandler() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET") //GET메소드일때 이 usersHandler가 불려라 정함
	mux.HandleFunc("/users", CreateUserHandler).Methods("POST")
	// mux.HandleFunc("/users/89", getUserInfo89Handler) 89가 아니라 아이디를 나타내는 {id:[0-9]+}문법으로 (고릴라)
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler)

	return mux
}
