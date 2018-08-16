package exwechat

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/wealthworks/go-tencent-api/client"
)

const (
	urlToken   = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	urlGetUser = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	urlAddUser = "https://qyapi.weixin.qq.com/cgi-bin/user/create"
	urlDelUser = "https://qyapi.weixin.qq.com/cgi-bin/user/delete"

	urlListDept = "https://qyapi.weixin.qq.com/cgi-bin/department/list"
	urlListUser = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist"
)

var (
	corpId, corpSecret string
)

func init() {
	corpId = os.Getenv("EXWECHAT_CORP_ID")
	corpSecret = os.Getenv("EXWECHAT_CORP_SECRET")
}

type API struct {
	c *client.Client
}

func NewAPI() *API {
	if corpId == "" || corpSecret == "" {
		log.Fatal("EXWECHAT_CORP_ID or EXWECHAT_CORP_SECRET are empty or not found")
	}
	c := client.NewClient(urlToken)
	c.SetContentType("application/json")
	c.SetCorp(corpId, corpSecret)
	return &API{c}
}

func (a *API) GetUser(userId string) (*User, error) {
	token, err := a.c.GetAuthToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s&userid=%s", urlGetUser, token, userId)

	body, err := a.c.Get(uri)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal(body, user)

	return user, err
}

func (a *API) AddUser(user *User) (err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", urlAddUser, token)
	var data []byte
	data, err = json.Marshal(user)
	if err != nil {
		return
	}

	_, err = a.c.Post(uri, data)
	return
}

func (a *API) DeleteUser(userId string) (err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&userid=%s", urlDelUser, token, userId)

	_, err = a.c.Get(uri)
	return
}

func (a *API) ListDepartment(id int) (data []Department, err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&id=%d", urlListDept, token, id)

	var body []byte
	body, err = a.c.Get(uri)
	if err != nil {
		return nil, err
	}

	var ret departmentResponse
	err = json.Unmarshal(body, &ret)

	if err == nil {
		data = ret.Department
	}

	return
}

func (a *API) ListUser(deptId int, incChild bool) (data []User, err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	fc := "0"
	if incChild {
		fc = "1"
	}
	uri := fmt.Sprintf("%s?access_token=%s&department_id=%d&fetch_child=%s", urlListUser, token, deptId, fc)

	var body []byte
	body, err = a.c.Get(uri)
	if err != nil {
		return nil, err
	}

	var ret usersResponse
	err = json.Unmarshal(body, &ret)

	if err == nil {
		data = ret.Users
	}

	return
}
