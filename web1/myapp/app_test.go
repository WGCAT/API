package myapp

import (
	"encoding/json"
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

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/89") //89번유저의 정보를 가져오고싶다를 http.Get으로 보냄
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:89") //89를 포함해야한다
	//사실 정해진 아이디가 아닌 유저아이디가 뒤에 붙었을때 맞게 와야함 URL.Path를 가지고

	resp, err := http.Get(ts.URL + "/users/56")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:56")
}
func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json", //포스트로 보냈고
		strings.NewReader(`{"first_name" : "sujin", "last_name":"lee", "email":"seed9878@gmail.com"}`)) // 유저데이터 보냈고
	assert.NoError(err)                               //에러가 없어야한다
	assert.Equal(http.StatusCreated, resp.StatusCode) // 만들어진거 확인했고
	//읽어서 만들어진거 확인해야지
	user := new(User)
	err := json.NewDecoder(resp.Body).Decode(user) //만약에 서버가 보낸 제이슨 코드에 문제가
	assert.NoError(err)
	assert.NotEqual(0, user.ID) //유저의 아이디가 0이 아니다

	id := user.ID //아이디를 가지고 다시 Get을 해본다
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.NotEqual(http.StatusOK, resp.StatusCode)
	user2 := new(User)
	err := json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.NotEqual(user.ID, user2.ID)
	assert.NotEqual(user.FirstName, user2.FirstName)
}
