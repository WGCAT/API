package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

var userMap map[int]*User
var lastID int //아이디 변수

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}
func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get UserInfo by /users/{id}")
}
func getUserInfoHandler(w http.ResponseWriter, r *http.Request) { //고정 89가 아니라 아이디를 나타내줘야하므로 mux.Vars사용한다
	vars := mux.Vars(r)                 //클라이언트가 리퀘스트한 아이디가 있는지 확인 후 있으면 유저정보 반환
	id, err := strconv.Atoi(vars["id"]) //스트링을 인티저로 바꿔줌
	if err != nil {                     //변환과정에서 에러가 있다 하면 리퀘스트가 잘못된거니까
		w.WriteHeader(http.StatusBadRequest) //배드리퀘스트 해주고 에러출력
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[id] //아이디에 해당하는 유저가 실제 맵에 있는지 확인
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return
	}

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
	lastID++
	user.ID = lastID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user
}

// NewHandler make a new myapp handler
func NewHandler() http.Handler {
	userMap = make(map[int]*User) //맵에 유저를 언제 등록할거냐 바로위 크리에이트 할때
	lastID = 0
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET") //GET메소드일때 이 usersHandler가 불려라 정함
	mux.HandleFunc("/users", CreateUserHandler).Methods("POST")
	// mux.HandleFunc("/users/89", getUserInfo89Handler) 89가 아니라 아이디를 나타내는 {id:[0-9]+}문법으로 (고릴라)
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler)

	return mux
}
