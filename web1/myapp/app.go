package myapp

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}
func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get UserInfo by /users/{id}")
}
func getUserInfoHandler(w http.ResponseWriter, r *http.Request) { //고정 89가 아니라 아이디를 나타내줘야하므로 mux.Vars사용한다
	vars := mux.Vars(r) //알아서 테스트 패싱해주ㅠㅁ
	fmt.Fprint(w, "User Id:", vars["id"])
}

// NewHandler make a new myapp handler
func NewHandler() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler)
	// mux.HandleFunc("/users/89", getUserInfo89Handler) 89가 아니라 아이디를 나타내는 {id:[0-9]+}문법으로 (고릴라)
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler)

	return mux
}
