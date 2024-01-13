package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

/*
[구현해볼것]
1. GET함수로 모든 User정보 Response
2. GET함수로 ID입력해서 해당ID값을 가진 User정보 Response
3. POST함수로 User정보 users map에 저장 후 Response
4. PUT함수로 User정보 전체 수정 후 Response
5. PUT함수로 User정보 일부 수정 후 Response
6. DELETE함수로 ID입력해서 해당ID값을 가진 User정보 삭제 후 Response
7. DELETE함수로 다수의 ID입력해서 해당ID값을 가진 User정보들 삭제 후 Response
*/

// var users = map[string]*User{} //맵선언
var users []User //배열선언

type User struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

func main() {
	//함수로 http요청 멀티플렉서 인스턴스 생성. 멀티플렉서(mux)로 url경로를 동시에 여러개 처리하고 여러 경로에 적절한 핸들러를 매칭시키기 위해 사용
	server := http.NewServeMux()
	server.HandleFunc("/user/", userHandler)
	server.HandleFunc("/users", GetUsersHandler)

	//서버의 포트를 지정하여 실제 웹 서버를 구동
	//첫번째는 포트 8091 에서 Request를 Listen 할 것을 지정
	//두번째는 어떤 ServeMux를 사용할 지를 지정 (nil인 경우 DefaultServeMux를 사용, 개발자가 별도로 ServeMux를 만들어 Routing 부분을 세밀하게 제어가능)
	http.ListenAndServe(":8091", server)

}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			pathParams := strings.Split(r.URL.Path, "/")
			userId := pathParams[len(pathParams)-1]

			//user변수 User(json형)자료형으로 선언
			var user User

			for _, u := range users { //forEach문 처럼 users배열의 값들을 루프돌림
				if u.ID == userId { //u의 ID값이 /주소값과 같을경우 user변수 값은 u값(주소값)
					user = u

				}
			}

			//json.NewEncoder(w) : 인코더를 생성. 생성된 인코더는 앞으로 입력할 데이터를 표준출력
			//json.Encode(user) : user변수값을 보내면 표준출력에는 json형으로 출력되는 인코딩 실행
			json.NewEncoder(w).Encode(user)

		}
	case "POST":
		{
			var user User

			//json.NewDecoder(r.Body) : 디코더 생성. 생성된 디코더는 앞으로 입력할 데이터를 표준입력
			//json.Decode(user) : 표준입력으로부터 json형 데이터가 들어오면 user변수값 데이터로 변경되는 디코딩 실행
			json.NewDecoder(r.Body).Decode(user)

			users = append(users, user) //배열의 add와 같은개념. users배열에 user값을 집어넣는다
			/*
			   s := []int{0, 1}
			   s = append(s, 2)        // 0,1,2
			   s = append(s, 3, 4, 5)  // 0,1,2,3,4,5
			*/

		}
	case "PUT":
		{
			var user User

			json.NewDecoder(r.Body).Decode(user)
			users = append(users, user)

		}
	case "DELETE":
		{
			pathParams := strings.Split(r.URL.Path, "/")
			userId := pathParams[len(pathParams)-1]

			var result []User
			// fmt.Println(result)
			// fmt.Println(users)
			for _, u := range users {

				if userId == u.ID {

					continue
				}
				result = append(result, u)
			}

		}

	}
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

/*맵 "map[Key타입] Value타입"
1. var map = map[string]string
{
   "ID":"1",
   "name":"aaa"
}
{
   "ID":"2",
   "name":"bbb"
}

2. var map = map[string]User
{
   "ID":"1" {
      "ID":"1",
     "name":"aaa"
   }
}
{
   "ID":"2" {
      "ID":"2",
     "name":"bbb"
   }
}

//배열 "var 변수명 [배열크기] 데이타타입"
1. var array []string
{
   "ID":"1",
   "name":"aaa"
}
{
   "ID":"2",
   "name":"bbb"
}
2. var array []User
{
   "ID":"1",
   "name":"aaa"
}
{
   "ID":"2",
   "name":"bbb"
}
*/
