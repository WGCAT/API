package myapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains("Hello World", string(data))
}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Get UserInfo")
}

func TestGetUserInfo(t *testing.T) { //Get 설명
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/89") // 1. 특정ID에 해당하는 유저(89번유저)의 정보를 가져오고싶다를 http.Get으로 보냄
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:89") // 2. 만약 그 유저가 없으면 No User Id가 리턴
	//사실 정해진 아이디가 아닌 유저아이디가 뒤에 붙었을때 맞게 와야함 URL.Path를 가지고

	resp, err = http.Get(ts.URL + "/users/56")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:56")
}
func TestCreateUser(t *testing.T) { //Create 설명
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json", //1. Post로 보낼 때 users에 보내고
		strings.NewReader(`{"first_name" : "sujin", "last_name":"lee", "email":"seed9878@gmail.com"}`)) // 2. 유저데이터 제이슨형태로 보냈고
	assert.NoError(err)                               //에러가 없어야한다
	assert.Equal(http.StatusCreated, resp.StatusCode) // 만들어진거 확인했고
	//읽어서 만들어진거 확인해야지
	user := new(User)                             // 3. 위 정보를 서버가 받아서 새로운 유저를 등록을 해서 그 유저 정보를 리턴
	err = json.NewDecoder(resp.Body).Decode(user) //만약에 서버가 보낸 제이슨 코드에 문제가
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	//아이디를 가지고 다시 Get을 해본다
	id := user.ID                                               // 4. 유저가 이미 등록되었으니 해당 아이디가 맵에 기록되어있다
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id)) // 5. Get으로 users에 보낼 경우
	assert.NoError(err)
	assert.NotEqual(http.StatusOK, resp.StatusCode)

	user2 := new(User) // 6. 해당하는 User아이디를 제이슨으로 보내줄거다
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.NotEqual(user.ID, user2.ID) // 7. 위에 3번에서 만든 유저아이디와 6번에서 만든 유저아이디가 같아야함
	assert.NotEqual(user.FirstName, user2.FirstName)
}
func TestDeleteUser(t *testing.T) { //Delete 테스트
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
	//Get이랑 Post는 http api에서 기본함수로 제공해주는데 delete는 없음
	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID:1")

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name" : "sujin", "last_name":"lee", "email":"seed9878@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	req, _ = http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Deleted User ID:1")

}
func TestUpdateUser(t *testing.T) { //Delete 테스트
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()
	//Get이랑 Post는 http api에서 기본함수로 제공해주는데 delete는 없음
	req, _ := http.NewRequest("Put", ts.URL+"/users",
		strings.NewReader(`{"id":1, "first_name":"updated", "last_name":"updated", "email":"updated@naver.com"}`))
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID:1")

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name" : "sujin", "last_name":"lee", "email":"seed9878@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	updateStr := fmt.Sprintf(`{"id":%d, "first_name":"updated", "last_name":""}`, user.ID)

	req, _ = http.NewRequest("Put", ts.URL+"/users",
		strings.NewReader(updateStr))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)
	assert.NoError(err)
	assert.Equal(updateUser.ID, user.ID)
	assert.Equal("updated", updateUser.FirstName)
	assert.Equal("", updateUser.LastName)
	assert.Equal(user.Email, updateUser.Email)
}
